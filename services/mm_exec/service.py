import traceback
from itertools import chain

import gen.proto.python.match_service.match_service_pb2 as ms_pb2
import gen.proto.python.match_service.match_service_pb2_grpc as ms_grpc_pb2
import models
import redis_repo


class MMExecService:
    def __init__(self, redis_ops: redis_repo.Ops, ms_client: ms_grpc_pb2.MatchServiceStub):
        self.redis_ops = redis_ops
        self.ms_client = ms_client

    def make_player_data(self, player: models.Player) -> ms_pb2.PlayerData:
        return ms_pb2.PlayerData(
            player_id=player.id,
            player_name=player.name,
            player_rating=player.displayed_rating,
            reg_id=player.reg_id,
        )

    def make_input_match(self, match: models.Match) -> ms_pb2.InputMatch:
        return ms_pb2.InputMatch(
            player1=self.make_player_data(match.player1),
            player2=self.make_player_data(match.player2),
            fleet=match.fleet
        )

    def make_request(self, matches: list[models.Match]) -> ms_pb2.StartMatchRequest:
        return ms_pb2.StartMatchRequest(
            matches=[
                self.make_input_match(match)
                for match in matches
            ]
        )

    def execute(self):
        self.redis_ops.reset()
        with self.redis_ops.lock():
            if not self.redis_ops.verify_time():
                return
            players = self.redis_ops.get_players()
            print("players in the pool:", [p.id for p in players])
            matches = self.redis_ops.gather_matches(players)
            print("matches: ", matches)
            to_return = set(str(p.id) for p in players) - set(
                chain.from_iterable((str(m.player1.id), str(m.player2.id)) for m in matches))
            ms_request = self.make_request(matches)
            try:
                response = self.ms_client.StartMatch(ms_request)
                print("ms response:", response)
                for inp_match, out_match in zip(matches, response.results):
                    if out_match.player1_fail_response == ms_pb2.PLAYER_FAIL_RESPONSE_REENTER:
                        to_return.add(str(inp_match.player1.id))
                    if out_match.player2_fail_response == ms_pb2.PLAYER_FAIL_RESPONSE_REENTER:
                        to_return.add(str(inp_match.player2.id))
                self.redis_ops.remove_players(p.id for p in players if str(p.id) not in to_return)
                self.redis_ops.return_players(len(players), to_return)
            except Exception:
                print("failed to request matches", traceback.format_exc())
            self.redis_ops.update_time()
            print("done")
