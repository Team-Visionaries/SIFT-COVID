from wtforms import Form, StringField

class SearchBar(Form):
    search = StringField('')