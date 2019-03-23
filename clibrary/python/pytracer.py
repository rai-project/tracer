#! /usr/bin/env python3

import sys
import os
import time
import logging
import ctypes
import inspect
import traceback
from distutils import sysconfig

logging.basicConfig()
logger = logging.getLogger("pytracer")
logger.setLevel(logging.WARN)


try:
    import threading

    def _settrace(func):
        threading.setprofile(func)
        sys.setprofile(func)

    def _unsettrace():
        sys.setprofile(None)
        threading.setprofile(None)


except ImportError:
    _settrace = sys.setprofile

    def _unsettrace():
        sys.setprofile(None)


class SpanStartFromContext_return(ctypes.Structure):
    _fields_ = [("span", ctypes.c_size_t), ("context", ctypes.c_size_t)]


RAITRACER_PATHS = [
    "librai_tracer.so",
    os.environ["GOPATH"]
    + "/src/github.com/rai-project/tracer/clibrary/dist/MacOSX-x86-64/librai_tracer.so",
]
NO_TRACE = 0
APPLICATION_TRACE = 1
MODEL_TRACE = 2
FRAMEWORK_TRACE = 3
LIBRARY_TRACE = 4
HARDWARE_TRACE = 5
FULL_TRACE = 6


class Tracer(object):
    def __init__(self, prog_argv):
        self.prog_argv = prog_argv
        self.libraitracer = None
        self.globalCtx = None
        self.globalSpanCtx = None
        for path in RAITRACER_PATHS:
            try:
                self.libraitracer = ctypes.cdll.LoadLibrary(path)
                self.libraitracer.TracerInit()
                logger.info("loaded RAI Tracer from {}".format(path))
                self.libraitracer.SpanStart.restype = ctypes.c_size_t
                self.libraitracer.SpanStart.argtypes = [ctypes.c_int, ctypes.c_char_p]
                self.libraitracer.SpanStartFromContext.restype = (
                    SpanStartFromContext_return
                )
                self.libraitracer.SpanStartFromContext.argtypes = [
                    ctypes.c_size_t,
                    ctypes.c_int,
                    ctypes.c_char_p,
                ]
                self.libraitracer.SpanAddTag.restype = None
                self.libraitracer.SpanAddTag.argtypes = [
                    ctypes.c_size_t,
                    ctypes.c_char_p,
                    ctypes.c_char_p,
                ]
                self.libraitracer.SpanFinish.restype = None
                self.libraitracer.SpanFinish.argtypes = [ctypes.c_size_t]

                self._spanStartFromContext(0, self.prog_argv[0])
                self.init_libpath()
                break
            except OSError as e:
                logger.debug("failed to load RAI Tracer from {}".format(path))
                self.libraitracer = None
        if not self.libraitracer:
            logger.error("couldn't load any of {}".format(RAITRACER_PATHS))

    def __del__(self):
        if self.globalSpanCtx is not None:
            # print("__del = ", self.globalSpanCtx.span)
            self._spanFinish(span_id=self.globalSpanCtx.span)
        if self.libraitracer is not None:
            self.libraitracer.TracerClose()

    def _spanStart(self, operationName):
        if self.globalSpanCtx is None:
            self.globalSpanCtx = self._spanStartFromContext(0, self.prog_argv[0])
        if self.libraitracer is not None:
            logger.debug("spanstart {}".format(operationName))
            # self.span_id = self.libraitracer.SpanStart(
            #     APPLICATION_TRACE, ctypes.c_char_p(str.encode(operationName))
            # )
            # return self.libraitracer.SpanStartFromContext(
            #     self.globalCtx,
            #     APPLICATION_TRACE,
            #     ctypes.c_char_p(str.encode(operationName)),
            # )
            res = self._spanStartFromContext(self.globalCtx, operationName)
            return res.span

    def _spanStartFromContext(self, ctx, operationName):
        if self.libraitracer is not None:
            logger.debug("spanstart {}".format(operationName))
            res = self.libraitracer.SpanStartFromContext(
                ctx, APPLICATION_TRACE, ctypes.c_char_p(str.encode(operationName))
            )
            if self.globalSpanCtx is None:
                self.globalSpanCtx = res
                self.globalCtx = self.globalSpanCtx.context
            return res

    def _spanFinish(self, span_id=None):
        if self.libraitracer is not None:
            if span_id is None:
                span_id = self.span_id
            logger.debug("spanfinish {}".format(span_id))
            self.libraitracer.SpanFinish(span_id)

    def _addTag(self, span_id, key, val):
        if self.libraitracer is not None:
            # print(span_id)
            # print(val)
            self.libraitracer.SpanAddTag(
                span_id,
                ctypes.c_char_p(str.encode(str(key))),
                ctypes.c_char_p(str.encode(str(val))),
            )

    def _add_function_full_name(self, span_id, function, module=None):
        if not module:
            module = inspect.getmodule(function)

        if module:
            self._addTag(span_id, "module", module.__name__)
        self._addTag(span_id, "function", function.__name__)

    def _add_full_name(self, span_id, frame, module=None):
        if not module:
            module = inspect.getmodule(frame)

        if module:
            self._addTag(span_id, "module", module.__name__)
            self._addTag(span_id, "module_path", module.__file__)

        try:
            class_name = frame.f_locals["self"].__class__.__name__
            self._addTag(span_id, "class", class_name)
        except KeyError:
            pass

        file_name = frame.f_globals.get("__file__", None)
        if file_name is not None:
            full_file_name = os.path.abspath(file_name)
            self._addTag(span_id, "full_file_name", full_file_name)

        line_number = frame.f_lineno
        self._addTag(span_id, "line_number", str(line_number))

        codename = frame.f_code.co_name
        # if codename != '<module>':  # top level usually
        #    name.append(codename)  # function or a method
        self._addTag(span_id, "code_name", str(codename))
        return

    def tracefunc(self, frame, event, arg, ranges=[[]], mode=[None]):

        # wait for the import to return
        if mode[0]:
            if event == "return":
                if frame == mode[0]:
                    logger.debug("import returned")
                    mode[0] = None
            return self.tracefunc

        if event == "call" or event == "c_call":
            if event == "call":
                # don't record call of _unsettrace (won't see exit)
                function_name = frame.f_code.co_name
                if function_name == "_unsettrace":
                    return self.tracefunc
            else:
                function_name = arg.__name__

            # skip builtins
            if inspect.isbuiltin(arg):
                return self.tracefunc

            # skip any imports
            if "importlib" in frame.f_code.co_filename:
                logger.debug("saw importlib..wait for return...")
                mode[0] = frame
                return self.tracefunc

            module = inspect.getmodule(frame)

            if module and self.is_module_stdlib(module.__file__):
                return self.tracefunc

            # if we have come across the init of a module, don't record ranges until it returns
            if module and function_name == "<module>":
                logger.debug("init of ", module, " wait for return...")
                mode[0] = frame
                return self.tracefunc

            if function_name == "<module>":
                return self.tracefunc

            # we may have defined the functions that are not part of a module, so we don't want to skip
            # if module is None:
            #     return tracefunc

            span_id = self._spanStart(function_name)
            if event == "call":
                self._add_full_name(span_id, frame, module=module)
            else:
                self._add_function_full_name(span_id, arg, module=module)
            # self._addTag(span_id, "frame", frame)
            self.span_id = span_id
            ranges[0].append(frame)
        elif event == "return" or event == "c_return":
            # don't record exit of _settrace (won't see call)
            # if frame.f_code.co_name == "_settrace":
            #     return tracefunc
            if ranges[0]:
                if ranges[0][-1] == frame:
                    self._spanFinish(span_id=self.span_id)
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

    def init_libpath(self):
        self.lib_path = sysconfig.get_python_lib()
        path = os.path.split(self.lib_path)
        if path[1] == "site-packages":
            self.lib_path = path[0]
        self.lib_path = self.lib_path.lower()

    def is_module_stdlib(self, file_name):
        """
        Returns True if the file_name is in the lib directory. Used to check
        if a function is in the standard library or not.
        """
        if "torch" in file_name.lower():
            return False
        return file_name.lower().startswith(self.lib_path)


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

    prog_argv = args.commands

    t = Tracer(prog_argv)

    t.run()

    del t


if __name__ == "__main__":
    main()
