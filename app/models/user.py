from sqlalchemy.orm import Mapped, mapped_column
from sqlalchemy import text
from uuid import uuid4
from sqlalchemy import String, Integer, Column
from uuid import UUID
from typing import Optional
from .base import Base


class User(Base):
    __tablename__ = "users"

    first_name: Mapped[str] = mapped_column(String(80))
    last_name: Mapped[str] = mapped_column(String(80))
    email: Mapped[str] = mapped_column(String(80))
    password: Mapped[str] = mapped_column(String)
    role = Column(String, default="teacher")