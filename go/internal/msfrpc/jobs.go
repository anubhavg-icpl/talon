package msfrpc

import "context"

// ListJobs mirrors job.list.
func (c *Client) ListJobs(ctx context.Context) (map[string]any, error) {
	return c.Call(ctx, "job.list")
}

// JobInfo mirrors job.info.
func (c *Client) JobInfo(ctx context.Context, id string) (map[string]any, error) {
	return c.Call(ctx, "job.info", id)
}

// StopJob mirrors job.stop.
func (c *Client) StopJob(ctx context.Context, id string) (map[string]any, error) {
	return c.Call(ctx, "job.stop", id)
}
