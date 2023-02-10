from models import User
from jinja2 import utils
from utils.constants import LOG
from flask_expects_json import expects_json
from pony.orm import db_session, commit, select
from middlewares.auth_validator import check_access_token
from flask import Blueprint, abort, g, make_response, jsonify
from utils.schema import email_regex, signup_schema, login_schema
import os, re, bcrypt, hashlib, shortuuid, jwt, datetime
from utils.errors import (
    INVALID_EMAIL,
    INVALID_USERNAME,
    USER_EXISTS,
    LOGIN_FAILED,
    JWT_GENERATION_ERROR,
)


auth = Blueprint("auth", __name__)


@auth.route("/signup", methods=["POST"])
@expects_json(signup_schema)
def signup():
    raw_user = g.data
    if not re.fullmatch(email_regex, raw_user.get("email")):
        abort(make_response(jsonify(INVALID_EMAIL), 400))

    validated_email = utils.escape(raw_user.get("email"))
    if validated_email != raw_user.get("email"):
        abort(make_response(jsonify(INVALID_EMAIL), 400))

    with db_session:
        existing_user = select(u for u in User if u.email == validated_email)
        if len(existing_user) > 0:
            abort(make_response(jsonify(USER_EXISTS), 400))

    validated_username = utils.escape(raw_user.get("username"))
    if validated_username != raw_user.get("username"):
        abort(make_response(jsonify(INVALID_USERNAME), 400))

    password_hash = bcrypt.hashpw(
        password=hashlib.sha256(raw_user.get("password").encode("utf-8")).digest(),
        salt=bcrypt.gensalt(),
    )

    with db_session:
        User(
            username=validated_username,
            email=validated_email,
            password=password_hash.decode("utf-8"),
            uid=shortuuid.ShortUUID().random(length=8),
        )
        commit()
    return {"message": "success"}


@auth.route("/login", methods=["POST"])
@expects_json(login_schema)
def login():
    raw_user = g.data
    with db_session:
        existing_user = select(u for u in User if u.email == raw_user.get("email"))[:]
        if len(existing_user) == 0:
            abort(make_response(jsonify(LOGIN_FAILED), 400))
        existing_user = existing_user[0]
        if not bcrypt.checkpw(
            hashlib.sha256(raw_user.get("password").encode("utf-8")).digest(),
            existing_user.password.encode("utf-8"),
        ):
            abort(make_response(jsonify(LOGIN_FAILED), 400))
        try:
            expire_at = (
                datetime.datetime.utcnow() + datetime.timedelta(days=2)
            ).timestamp()
            payload = {"exp": expire_at, "uid": existing_user.uid}
            token = jwt.encode(payload, os.getenv("JWT_SECRET"), algorithm="HS256")
            resp = make_response(jsonify({"user": existing_user.uid}), 200)
            resp.set_cookie("access_token", token, max_age=expire_at)
            return resp
        except Exception as e:
            LOG.error(f"Failed to generate jwt token: {e}")
            abort(make_response(jsonify(JWT_GENERATION_ERROR), 400))


@auth.route("/logout", methods=["GET"])
def logout():
    resp = make_response({"message": "success"}, 200)
    resp.set_cookie("access_token", "", max_age=0)


@auth.route("/refresh", methods=["GET"])
def refresh():
    check_access_token()
    expire_at = (datetime.datetime.utcnow() + datetime.timedelta(days=2)).timestamp()
    payload = {"exp": expire_at, "uid": g.uid}
    token = jwt.encode(payload, os.getenv("JWT_SECRET"), algorithm="HS256")
    resp = make_response(jsonify({"user": g.uid}), 200)
    resp.set_cookie("access_token", token, max_age=expire_at)
    return resp
