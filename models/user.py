from models.db import db
from pony.orm import Required


class User(db.Entity):
    email = Required(str)
    username = Required(str)
    password = Required(str)
    uid = Required(str)
