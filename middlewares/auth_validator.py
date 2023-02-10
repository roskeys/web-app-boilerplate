import os, jwt, time
from flask_http_middleware import BaseHTTPMiddleware
from utils.errors import NO_COOKIE_FOUND, ACCESS_TOKEN_EXPIRED
from flask import request, abort, make_response, jsonify, g, Request, Response


def check_access_token():
    if not request.cookies:
        abort(make_response(jsonify(NO_COOKIE_FOUND), 401))
    access_token = request.cookies.get("access_token")
    payload = jwt.decode(
        access_token,
        os.getenv("JWT_SECRET").encode("utf-8"),
        algorithms="HS256",
    )
    if time.time() > payload.get("exp"):
        abort(make_response(jsonify(ACCESS_TOKEN_EXPIRED), 401))
    g.uid = payload.get("uid")


class AccessCheckerMiddleware(BaseHTTPMiddleware):
    def __init__(self):
        super().__init__()

    def dispatch(self, request: Request, call_next: callable) -> Response:
        if request.path.startswith("/auth/"):
            return call_next(request)
        check_access_token()
        return call_next(request)
