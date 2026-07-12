package strike

import "context"

// ListSessions mirrors session.list: the response map's own keys ARE the
// session IDs (no wrapper key), each value a session info map including a
// "type" field ("shell" or "meterpreter").
//
// Session IDs arrive as msgpack integers; Client.doCall normalizes them to
// string keys so callers can look up strconv.Itoa(id) consistently.
func (c *Client) ListSessions(ctx context.Context) (map[string]any, error) {
	return c.Call(ctx, "session.list")
}

func readMethod(sessionType string) string {
	if sessionType == "meterpreter" {
		return "session.meterpreter_read"
	}
	return "session.shell_read"
}

func writeMethod(sessionType string) string {
	if sessionType == "meterpreter" {
		return "session.meterpreter_write"
	}
	return "session.shell_write"
}

// ReadSession dispatches to session.shell_read or session.meterpreter_read
// based on sessionType.
func (c *Client) ReadSession(ctx context.Context, id, sessionType string) (map[string]any, error) {
	return c.Call(ctx, readMethod(sessionType), id)
}

// WriteSession dispatches to session.shell_write or session.meterpreter_write.
func (c *Client) WriteSession(ctx context.Context, id, sessionType, data string) (map[string]any, error) {
	return c.Call(ctx, writeMethod(sessionType), id, data)
}

// StopSession mirrors session.stop.
func (c *Client) StopSession(ctx context.Context, id string) (map[string]any, error) {
	return c.Call(ctx, "session.stop", id)
}
