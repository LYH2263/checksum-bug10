package checksum

import "context"

func DigestStream(ctx context.Context, opts Options, data []byte) (*Manifest, error) {
	p := New(opts)
	if len(data) == 0 {
		return nil, ErrEmpty
	}
	if err := p.Feed(ctx, data); err != nil {
		return nil, err
	}
	return p.Finish()
}

func ServeDigest(ctx context.Context, p *Pipeline, body []byte) (*Manifest, error) {
	return DigestStream(ctx, p.opts, body)
}
