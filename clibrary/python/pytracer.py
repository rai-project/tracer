#! /usr/bin/env python3

import sys
import os
import time
import logging
import ctypes
import inspect
import traceback

logging.basicConfig()
logger = logging.getLogger("pytracer")
logger.setLevel(logging.DEBUG)

DEPTH_LIMIT = None


try:
    import threading
except ImportError:
    _settrace = sys.setprofile

    def _unsettrace():
        sys.setprofile(None)


else:

    def _settrace(func):
        threading.setprofile(func)
        sys.setprofile(func)

    def _unsettrace():
        sys.setprofile(None)
        threading.setprofile(None)


NO_TRACE = 0
APPLICATION_TRACE = 1
MODEL_TRACE = 2
FRAMEWORK_TRACE = 3
LIBRARY_TRACE = 4
HARDWARE_TRACE = 5
FULL_TRACE = 6


def function_full_name(function, module=None):
    name = []
    if not module:
        module = inspect.getmodule(function)
    # print(inspect.getmembers(module))

    if module:
        name.append(module.__name__)

    name.append(function.__name__)
    return name


def full_name(frame, module=None):
    name = []
    if not module:
        module = inspect.getmodule(frame)
    # print(inspect.getmembers(module))

    if module:
        name.append(module.__name__)
    if "self" in frame.f_locals:
        # I don't know any way to detect call from the object method
        # XXX: there seems to be no way to detect static method call - it will
        #      be just a function call
        try:
            class_name = frame.f_locals["self"].__class__.__name__
        except KeyError:
            class_name = None
        if class_name:
            name.append(class_name)
    codename = frame.f_code.co_name
    # if codename != '<module>':  # top level usually
    #    name.append(codename)  # function or a method
    name.append(codename)
    return name


class Tracer(object):
    def __init__(self, prog_argv, depth_limit=None):
        self.prog_argv = prog_argv
        self.depth_limit = depth_limit
        self.libraitracer = None
        self.libcudart = None
        # load the rai tracer library
        RAITRACER_PATHS = [
            "librai_tracer.so",
            os.environ["GOPATH"]
            + "/src/github.com/rai-project/tracer/clibrary/dist/MacOSX-x86-64/librai_tracer.so",
            "/usr/local/lib/librai_tracer.so",
        ]
        for path in RAITRACER_PATHS:
            try:
                self.libraitracer = ctypes.cdll.LoadLibrary(path)
            except OSError as e:
                logger.debug("failed to load RAI Tracer from {}".format(path))
                self.libraitracer = None
            else:
                logger.info("loaded RAI Tracer from {}".format(path))
                break
        if not self.libraitracer:
            logger.error("couldn't load any of {}".format(RAITRACER_PATHS))

    def _spanStart(self, operationName):
        if self.libraitracer:
            logger.debug("spanstart {}".format(operationName))
            # self.libraitracer.SpanStart(
            #     ctypes.c_int(p0), ctypes.c_char_p(str.encode(p1))
            # )

    def _spanFinish(self, operationName):
        if self.libraitracer:
            logger.debug("spansfinish {}".format(operationName))
            # self.libraitracer.SpanFinish(ctypes.c_void_p(p0))

    def tracefunc(self, frame, event, arg, ranges=[[]], mode=[None]):

        # wait for the import to return
        if mode[0]:
            if event == "return":
                if frame == mode[0]:
                    logger.debug("import returned")
                    mode[0] = None
            return self.tracefunc

        if event == "call" or event == "c_call":
            if self.depth_limit:
                if len(ranges[0]) > self.depth_limit:
                    return self.tracefunc

            if event == "call":
                # don't record call of _unsettrace (won't see exit)
                function_name = frame.f_code.co_name
                if function_name == "_unsettrace":
                    return self.tracefunc
            else:
                # skip builtins
                if inspect.isbuiltin(arg):
                    return self.tracefunc
                function_name = arg.__name__

            # skip any imports
            if "importlib" in frame.f_code.co_filename:
                logger.debug("saw importlib..wait for return...")
                mode[0] = frame
                return self.tracefunc

            module = inspect.getmodule(frame)

            # if we have come across the init of a module, don't record ranges until it returns
            if module and function_name == "<module>":
                logger.debug("init of ", module, " wait for return...")
                mode[0] = frame
                return self.tracefunc

            # if function_name == "<module>":
            #     return tracefunc

            # we may have defined the functions that are not part of a module, so we don't want to skip
            # if module is None:
            #     return tracefunc

            if event == "call":
                name = []
                # name += [str(frame.f_code.co_filename)]
                name += full_name(frame, module=module)
            else:
                name = function_full_name(arg, module=module)
            # filename, lineno, function_name, code_context, index = inspect.getframeinfo(frame)
            range_name = ".".join(name)
            ranges[0].append(frame)
            self._spanStart(range_name)
        # elif event == "c_return":
        #     # arg is the c function object
        #     frame_depth = depth[0]
        #     depth[0] -= 1
        #     if DEPTH_LIMIT and depth[0] > DEPTH_LIMIT:
        #         return tracefunc
        elif event == "return" or event == "c_return":
            # don't record exit of _settrace (won't see call)
            # if frame.f_code.co_name == "_settrace":
            #     return tracefunc
            if ranges[0]:
                if ranges[0][-1] == frame:
                    self._spanFinish(".".join(full_name(frame)))
                    ranges[0] = ranges[0][:-1]
                    # name = full_name(frame)

        return self.tracefunc

    def runctx(self, cmd, globals=None, locals=None):
        if globals is None:
            globals = {}
        if locals is None:
            locals = {}
        _settrace(self.tracefunc)
        try:
            exec(cmd, globals, locals)
        finally:
            _unsettrace()

    def run(self):
        logger.debug("Tracing python argv[:] {}".format(self.prog_argv))
        sys.argv = self.prog_argv
        progname = self.prog_argv[0]
        sys.path[0] = os.path.split(progname)[0]
        try:
            with open(progname) as fp:
                code = compile(fp.read(), progname, "exec")
                # try to emulate __main__ namespace as much as possible
                globs = {
                    "__file__": progname,
                    "__name__": "__main__",
                    "__package__": None,
                    "__cached__": None,
                }
                self.runctx(code, globs, globs)
        except IOError as err:
            logger.critical(
                "Cannot add rai tracer to python file %r because: %s"
                % (sys.argv[0], err)
            )
            sys.exit(1)
        except SystemExit:
            pass


def main():

    import argparse

    parser = argparse.ArgumentParser(
        description="Add RAI Tracer ranges to python functions"
    )
    parser.add_argument(
        "--depth", type=int, help="only push ranges to this stack depth"
    )
    parser.add_argument("--debug", action="store_true", help="print debug messages")
    parser.add_argument("--verbose", action="store_true", help="print verbose messages")
    parser.add_argument("commands", nargs="+", help="commands help")

    args = parser.parse_args()
    if args.debug:
        logger.setLevel(logging.DEBUG)
    if args.verbose:
        logger.setLevel(logging.INFO)
    if args.depth:
        if args.depth < 0:
            logger.critical("trace depth must be >=0")
            sys.exit(1)
        else:
            DEPTH_LIMIT = args.depth

    prog_argv = args.commands

    t = Tracer(prog_argv, depth_limit=args.depth)

    t.run()


if __name__ == "__main__":
    main()
