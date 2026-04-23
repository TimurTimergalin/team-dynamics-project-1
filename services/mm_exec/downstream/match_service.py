from abc import ABC, abstractmethod

import grpc

import gen.proto.python.match_service.match_service_pb2 as ms_pb2
import gen.proto.python.match_service.match_service_pb2_grpc as ms_grpc_pb2


class MatchServiceClientFactory(ABC):
    @abstractmethod
    def start_match(self, request: ms_pb2.StartMatchRequest) -> ms_pb2.StartMatchResponse:
        ...


class GrpcMatchServiceClientFactory(MatchServiceClientFactory):
    def __init__(self, address: str):
        self._address = address

    def start_match(self, request: ms_pb2.StartMatchRequest) -> ms_pb2.StartMatchResponse:
        with grpc.insecure_channel(self._address) as channel:
            return ms_grpc_pb2.MatchServiceStub(channel).StartMatch(request)
