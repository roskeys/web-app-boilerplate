from models.db import db
from models.user import User

db.bind("sqlite", "test.sqlite", create_db=True)
db.generate_mapping(create_tables=True)
