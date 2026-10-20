"""Forward tool events independently of cmd/go's buffered tool output."""

import threading


class ToolEventStream:
    def __init__(self, path, output):
        self.path = path
        self.output = output
        self.stop = threading.Event()
        self.thread = threading.Thread(target=self.stream, daemon=True)

    def __enter__(self):
        self.path.touch(exist_ok=False)
        self.thread.start()
        return self

    def stream(self):
        pending = ""
        with self.path.open() as source:
            while True:
                pending += source.read()
                lines = pending.split("\n")
                pending = lines.pop()
                for line in lines:
                    self.output.write("[GO TOOL] " + line + "\n")
                    self.output.flush()
                if self.stop.is_set():
                    # The producer has exited before __exit__ signals us. Read
                    # once more in case it appended during the preceding read.
                    pending += source.read()
                    for line in pending.splitlines():
                        self.output.write("[GO TOOL] " + line + "\n")
                        self.output.flush()
                    return
                self.stop.wait(0.1)

    def __exit__(self, *_):
        self.stop.set()
        self.thread.join()
