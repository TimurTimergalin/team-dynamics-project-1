import pytest
import services.mm_exec.config as config_py
import fakeredis


@pytest.fixture
def _start_match():
    def base_impl(*args, **kwargs):
        raise NotImplementedError("Not implemented")

    class StartMatchMockContainer:
        pass

    res = StartMatchMockContainer()
    res.start_match = base_impl
    return res


@pytest.fixture
def ms_client(_start_match):
    class MockClient:
        def StartMatch(self, *args, **kwargs):
            return _start_match.start_match(self, *args, **kwargs)

    return MockClient()


@pytest.fixture
def mock_start_match(_start_match):
    def dec(f):
        _start_match.start_match = f

    return dec


@pytest.fixture
def cfg():
    return config_py.MMExecConfig(10_000, 10_000, 5_000, "", "")


@pytest.fixture
def rdb():
    return fakeredis.FakeStrictRedis(decode_responses=True)
