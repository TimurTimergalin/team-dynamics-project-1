import time
import sys
import traceback
import datetime


class Scheduler:
    def __init__(self, job, period_millis):
        self.job = job
        self.period_millis = period_millis

    def run(self):
        while True:
            start = time.time_ns() // 1_000_000
            try:
                self.job()
            except Exception:
                print("Job failed:", traceback.format_exc(), file=sys.stderr, sep='\n')
            else:
                print("Job executed at:", datetime.datetime.now(datetime.timezone.utc))
            end = time.time_ns() // 1_000_000
            time.sleep(max(self.period_millis - end + start, 0) / 1_000)
