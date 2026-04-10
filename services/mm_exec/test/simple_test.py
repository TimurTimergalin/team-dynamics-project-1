from typing import NamedTuple

from google.protobuf import json_format

import services.mm_exec.redis_repo as redis_repo
import services.mm_exec.service as service_py
import gen.proto.python.match_service.match_service_pb2 as ms_pb2


class PlayerValues(NamedTuple):
    rating: float
    fleet: str
    name: str
    displayed_rating: int
    reg_id: str


def test_simple(ms_client, mock_start_match, cfg, rdb):
    @mock_start_match
    def _mock_start_match(client, request, *args, **kwargs):
        r_json = json_format.MessageToDict(request)
        assert r_json == {
            "matches": [
                {
                    "player1": {
                        "playerId": "3",
                        "playerName": "p3",
                        "playerRating": "828",
                        "regId": "rg3"
                    },
                    "player2": {
                        "playerId": "4",
                        "playerName": "p4",
                        "playerRating": "1129",
                        "regId": "rg4"
                    },
                    "fleet": "fleet1"
                },
                {
                    "player1": {
                        "playerId": "1",
                        "playerName": "p1",
                        "playerRating": "1200",
                        "regId": "rg1"
                    },
                    "player2": {
                        "playerId": "2",
                        "playerName": "p2",
                        "playerRating": "1400",
                        "regId": "rg2"
                    },
                    "fleet": "fleet1"
                },
                {
                    "player1": {
                        "playerId": "7",
                        "playerName": "p7",
                        "playerRating": "1200",
                        "regId": "rg7"
                    },
                    "player2": {
                        "playerId": "6",
                        "playerName": "p6",
                        "playerRating": "1400",
                        "regId": "rg6"
                    },
                    "fleet": "fleet2"
                }
            ]
        }
        resp = ms_pb2.StartMatchResponse(
            results=[
                {
                    "match_id": "1"
                },
                {
                    "player1_fail_response": "PLAYER_FAIL_RESPONSE_REENTER",
                    "player2_fail_response": "PLAYER_FAIL_RESPONSE_REMOVE",
                },
                {
                    "match_id": "2"
                }
            ]
        )
        return resp

    ser = service_py.MMExecService(
        redis_repo.Ops(rdb, cfg),
        ms_client
    )

    p_keys = [
        redis_repo.PlayerKeySet(1),
        redis_repo.PlayerKeySet(2),
        redis_repo.PlayerKeySet(3),
        redis_repo.PlayerKeySet(4),
        redis_repo.PlayerKeySet(5),
        redis_repo.PlayerKeySet(6),
        redis_repo.PlayerKeySet(7),
    ]
    p_vals = [
        PlayerValues(1200, "fleet1", "p1", 1200, "rg1"),
        PlayerValues(1400.285, "fleet1", "p2", 1400, "rg2"),
        PlayerValues(827.785, "fleet1", "p3", 828, "rg3"),
        PlayerValues(1129, "fleet1", "p4", 1129, "rg4"),
        PlayerValues(2000, "fleet1", "p5", 2000, "rg5"),
        PlayerValues(1400.285, "fleet2", "p6", 1400, "rg6"),
        PlayerValues(1200.285, "fleet2", "p7", 1200, "rg7"),
    ]

    for ks, vs in zip(p_keys, p_vals):
        rdb.rpush("mmpool", ks.id)
        for k, v in zip(ks.keys(), vs):
            rdb.set(k, v)

    ser.execute()
    assert set(rdb.lrange("mmpool", 0, -1)) == {"5", "1"}
    expected_presence = [
        True,
        False,
        False,
        False,
        True,
        False,
        False,
    ]

    for ks, pr in zip(p_keys, expected_presence):
        for k in ks.keys():
            assert (rdb.get(k) is not None) == pr
