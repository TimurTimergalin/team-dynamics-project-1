import pytest
import services.mm_exec.config as config_py
import services.mm_exec.downstream.match_service as ms_downstream
import fakeredis


@pytest.fixture
def mock_start_match():
    class MockMatchServiceFactory(ms_downstream.MatchServiceClientFactory):
        def __init__(self):
            self._impl = None

        def start_match(self, request):
            if self._impl is None:
                raise NotImplementedError("start_match not configured")
            return self._impl(request)

        def __call__(self, f):
            self._impl = f

    return MockMatchServiceFactory()


@pytest.fixture
def ms_factory(mock_start_match):
    return mock_start_match


@pytest.fixture
def cfg():
    return config_py.MMExecConfig(10_000, 10_000, 5_000, "")


@pytest.fixture
def rdb():
    return fakeredis.FakeStrictRedis(decode_responses=True)
