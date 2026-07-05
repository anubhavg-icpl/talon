from celery import Celery
import logging
import asyncio
import datetime
import sys
import os

# Add the project root to sys.path so 'pentest_logic' can be imported by Celery workers
project_root = os.path.dirname(os.path.abspath(__file__))
if project_root not in sys.path:
    sys.path.append(project_root)

# Initialize Agent-Specific Celery App
celery_app = Celery(
    "pentesting_worker",
    broker="amqp://guest:guest@localhost:5672/"
)

celery_app.conf.update(
    task_serializer="json",
    accept_content=["json"],
    result_serializer="json",
    task_ignore_result=True,
)

logger = logging.getLogger(__name__)

@celery_app.task(name="execute_agent_task")
def execute_pentest_task(run_id, project_id, agent_name, agent_inputs, **kwargs):
    """
    Isolated executor for Pentesting. 
    Receives JSON, returns JSON. Does NOT import app.models.
    """
    import sys
    import os
    # Ensure project root is in sys.path inside the worker process context
    project_root = os.path.dirname(os.path.abspath(__file__))
    if project_root not in sys.path:
        sys.path.append(project_root)
        
    from pentest_core.integration import run_hexstrike_agent
    
    logger.info(f"[PentestWorker] Received execution request for run {run_id}")
    
    try:
        # Execute the agent logic
        # Note: integration.py handles publishing to agent.pentest.output
        loop = asyncio.get_event_loop()
        if loop.is_closed():
            loop = asyncio.new_event_loop()
            asyncio.set_event_loop(loop)
            
        result = loop.run_until_complete(run_hexstrike_agent(
            target_ip=agent_inputs.get("target_ip"),
            cve_id=agent_inputs.get("cve_id"),
            lhost=agent_inputs.get("lhost", "0.0.0.0"),
            lport=int(agent_inputs.get("lport", 4444)),
            description=agent_inputs.get("description", ""),
            run_id=run_id,
            project_id=project_id
        ))

        # Report results back to the platform via the output queue
        if run_id and project_id:
            try:
                final_message = result["messages"][-1].content if result.get("messages") else "No output"
                
                completion_payload = {
                    "event_type": "RUN_COMPLETED",
                    "run_id": str(run_id),
                    "project_id": str(project_id),
                    "agent_name": agent_name,
                    "overall_status": "completed",
                    "result": {
                        "summary": final_message,
                        "raw_output": str(result)
                    },
                    "timestamp": datetime.datetime.utcnow().isoformat()
                }

                # Use the local celery_app instance to dispatch to the platform
                celery_app.send_task(
                    "get_pentest_ouput",
                    kwargs=completion_payload,
                    queue="agent.pentest.output"
                )
                logger.info(f"[PentestWorker] Reported completion for run {run_id}")
            except Exception as re:
                logger.error(f"[PentestWorker] Failed to report results: {re}")
        
        return {"status": "success", "run_id": run_id}
    except Exception as e:
        logger.error(f"[PentestWorker] Execution failed: {e}")
        # In case of failure, we should ideally notify the output queue as well
        # but integration.py usually does its own error handling.
        return {"status": "error", "error": str(e), "run_id": run_id}
