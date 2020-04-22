import os
import requests
import time


RATING_API = os.getenv("RATING_API", "http://rating-api:8000/graphql")
RATING_AGENT_JANUS = os.getenv(
    "RATING_AGENT_JANUS", "http://rating-agent-janus:8080/api/v1/janus-gateway"
)


def test_janus():
    query = """
mutation {
    a1:upsertAccount(
        name: "Fabio",
        account_tag: "1000",
        type: PREPAID,
        balance: 1000000,
        active: true,
        max_concurrent_transactions: 100
    ) {
        id
    }
    a2:upsertAccount(
        name: "Alex",
        account_tag: "1001",
        type: PREPAID,
        balance: 1000000,
        active: true,
        max_concurrent_transactions: 100
    ) {
        id
    }
    upsertCarrier(
        carrier_tag: "carrier",
        active: true,
        protocol: UDP,
        host: "carrier",
        port: 5060
    ) {
        id
    }
    upsertPricelist(
        pricelist_tag: "pricelist",
        currency:EUR
    ) {
        id
    }
    upsertPricelistRate(
        carrier_tag: "carrier",
        pricelist_tag: "pricelist",
        prefix: "sip:1001",
        active: true,
        description: "pricelist rate",
        rate: 1,
        rate_increment: 1,
        connect_fee: 0,
        interval_start: 0
    ) {
        id
    }
}"""
    r = requests.post(RATING_API, json={"query": query})
    assert r.status_code == 200
    #
    data = [
        {
            "type": 64,
            "timestamp": 1587492912000000,
            "session_id": 3446611678650423,
            "handle_id": 97153772170402,
            "opaque_id": "siptest-fHYbsADjXHzF",
            "event": {
                "plugin": "janus.plugin.sip",
                "data": {
                    "event": "sip-in",
                    "sip": """INVITE sip:service@192.168.192.5:5060 SIP/2.0
Via: SIP/2.0/UDP 192.168.192.2:5060;branch=z9hG4bK-18-1-0
From: sipp <sip:1000@192.168.192.2:5060>;tag=1
To: sut <sip:1001@anotherdomain.com:5060>
Call-ID: 1-18@192.168.192.2
CSeq: 1 INVITE
Contact: sip:1000@192.168.192.2:5060
Max-Forwards: 70
Subject: Test
Content-Type: application/sdp
Content-Length:   137

v=0
o=user1 53655765 2353687637 IN IP4 192.168.192.2
s=-
c=IN IP4 192.168.192.2
t=0 0
m=audio 6000 RTP/AVP 0
a=rtpmap:0 PCMU/8000

""".replace(
                        "\n", "\r\n"
                    ),
                },
            },
        },
        {
            "type": 64,
            "timestamp": 1587492913000000,
            "session_id": 3446611678650423,
            "handle_id": 97153772170402,
            "opaque_id": "siptest-fHYbsADjXHzF",
            "event": {
                "plugin": "janus.plugin.sip",
                "data": {
                    "event": "sip-in",
                    "sip": """ACK sip:192.168.192.3:5060;transport=UDP SIP/2.0
Via: SIP/2.0/UDP 192.168.192.2:5060;branch=z9hG4bK-18-1-0
From: sipp <sip:1000@192.168.192.2:5060>;tag=1
To: sut <sip:1001@anotherdomain.com:5060>;tag=1
Route: <sip:192.168.192.5;lr;did=a11.a121>
Call-ID: 1-18@192.168.192.2
CSeq: 1 ACK
Contact: <sip:1000@192.168.192.2:5060;transport=UDP>
Max-Forwards: 70
Subject: Test
Content-Length: 0

""".replace(
                        "\n", "\r\n"
                    ),
                },
            },
        },
        {
            "type": 64,
            "timestamp": 1587492914000000,
            "session_id": 3446611678650423,
            "handle_id": 97153772170402,
            "opaque_id": "siptest-fHYbsADjXHzF",
            "event": {
                "plugin": "janus.plugin.sip",
                "data": {
                    "event": "sip-in",
                    "sip": """BYE sip:192.168.192.3:5060;transport=UDP SIP/2.0
Via: SIP/2.0/UDP 192.168.192.2:5060;branch=z9hG4bK-18-1-6
Route: <sip:192.168.192.5;lr;did=a11.a121>
From: sipp <sip:1000@192.168.192.2:5060>;tag=1
To: sut <sip:1001@anotherdomain.com:5060>;tag=1
Call-ID: 1-18@192.168.192.2
CSeq: 2 BYE
Contact: <sip:192.168.192.2:5060;transport=UDP>
Max-Forwards: 30
Content-Length: 0

""".replace(
                        "\n", "\r\n"
                    ),
                },
            },
        },
    ]
    for item in data:
        # sleep 500 ms
        time.sleep(0.5)
        # send the message
        r = requests.post(RATING_AGENT_JANUS, json=[item])
        assert r.status_code == 201
    # wait 2s for consolidation
    time.sleep(2.0)
    # verify the transactions
    query = """
query {
  allTransactions {
    transaction_tag
    account_tag
    inbound
    primary
    source
    destination
    fee
    duration
    timestamp_auth
    timestamp_begin
    timestamp_end
  }
}"""
    r = requests.post(RATING_API, json={"query": query})
    assert r.status_code == 200
    response = r.json()
    assert response["data"]["allTransactions"] == [
        {
            "transaction_tag": "1-18@192.168.192.2",
            "account_tag": "1000",
            "inbound": False,
            "primary": True,
            "source": "sip:1000@192.168.192.2",
            "destination": "sip:1001@anotherdomain.com",
            "fee": 1,
            "duration": 1,
            "timestamp_auth": None,
            "timestamp_begin": "2020-04-21T18:15:13",
            "timestamp_end": "2020-04-21T18:15:14",
        },
        {
            "transaction_tag": "1-18@192.168.192.2",
            "account_tag": "1001",
            "inbound": True,
            "primary": True,
            "source": "sip:1000@192.168.192.2",
            "destination": "sip:1001@anotherdomain.com",
            "fee": 0,
            "duration": 1,
            "timestamp_auth": None,
            "timestamp_begin": "2020-04-21T18:15:13",
            "timestamp_end": "2020-04-21T18:15:14",
        },
    ]
