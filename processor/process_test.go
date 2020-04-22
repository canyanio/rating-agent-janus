package processor

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/canyanio/rating-agent-janus/client/rabbitmq/mock"
	"github.com/canyanio/rating-agent-janus/model"
	"github.com/canyanio/rating-agent-janus/state"
)

func TestNewJanusProcessor(t *testing.T) {
	srv := NewJanusProcessor(nil, nil)
	assert.NotNil(t, srv)
}

func TestProcessSIP(t *testing.T) {
	ctx := context.Background()
	payload := []byte(`{"emitter":"janus-gateway","type":64,"timestamp":1587492912377714,"session_id":3446611678650423,"handle_id":97153772170402,"opaque_id":"siptest-fHYbsADjXHzF","event":{"plugin":"janus.plugin.sip","data":{"event":"sip-out","sip":"INVITE sip:3292166164@sip.messagenet.it SIP/2.0\r\nVia: SIP/2.0/TCP 192.168.48.2:50310;rport;branch=z9hG4bKyrByZ50mrv41S\r\nMax-Forwards: 70\r\nFrom: \"Fabio\" <sip:5313891@sip.messagenet.it>;tag=vp7H4t214XcaD\r\nTo: <sip:3292166164@sip.messagenet.it>\r\nCall-ID: 9vyJcXyHZzsRwpC5iUzahEc\r\nCSeq: 949120423 INVITE\r\nContact: Fabio <sip:5313891@192.168.48.2:50310>\r\nUser-Agent: Janus WebRTC Server SIP Plugin 0.0.8\r\nAllow: INVITE, ACK, BYE, CANCEL, OPTIONS, UPDATE, REFER, MESSAGE, INFO, NOTIFY\r\nSupported: replaces\r\nContent-Type: application/sdp\r\nContent-Disposition: session\r\nContent-Length: 978\r\n\r\nv=0\r\no=- 954527097879086087 4075993287627398736 IN IP4 192.168.48.2\r\ns=-\r\nt=0 0\r\nm=audio 28518 RTP/AVP 111 103 104 9 0 8 106 105 13 110 112 113 126\r\nc=IN IP4 192.168.48.2\r\na=rtpmap:111 opus/48000/2\r\na=fmtp:111 minptime=10;useinbandfec=1\r\na=rtpmap:103 ISAC/16000\r\na=rtpmap:104 ISAC/32000\r\na=rtpmap:9 G722/8000\r\na=rtpmap:0 PCMU/8000\r\na=rtpmap:8 PCMA/8000\r\na=rtpmap:106 CN/32000\r\na=rtpmap:105 CN/16000\r\na=rtpmap:13 CN/8000\r\na=rtpmap:110 telephone-event/48000\r\na=rtpmap:112 telephone-event/32000\r\na=rtpmap:113 telephone-event/16000\r\na=rtpmap:126 telephone-event/8000\r\na=extmap:1 urn:ietf:params:rtp-hdrext:ssrc-audio-level\r\na=extmap:2 http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time\r\na=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\na=extmap:4 urn:ietf:params:rtp-hdrext:sdes:mid\r\na=extmap:5 urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id\r\na=extmap:6 urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id\r\na=rtcp-fb:111 transport-cc\r\n"}}}`)

	event := &model.Event{}
	json.Unmarshal(payload, event)

	stateManager := state.NewMemoryManager()
	client := &mock.Client{}

	srv := NewJanusProcessor(stateManager, client)
	err := srv.Process(ctx, event)
	assert.Nil(t, err)
}
