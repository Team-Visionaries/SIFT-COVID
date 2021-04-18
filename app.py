import json
import newspaper #https://newspaper.readthedocs.io/en/latest/
from newspaper import Article 
from search import SearchBar
from bs4 import BeautifulSoup # https://www.crummy.com/software/BeautifulSoup/bs4/doc/
from flask import Flask, flash, render_template, request, redirect

app = Flask(__name__)

@app.route('/tool')
def tool(search):
    url = search.data['search']

    if url == '' or not url.startswith('http'):
            flash('Please enter a URL!')
            return redirect('/search')
    else:
        article = Article(url)
        article.download()
        article.parse()
        article.nlp()

        results = {
            'url': url,
            'title': article.title,
            'date:': article.publish_date,
            'keywords': article.keywords,
            'summary': article.summary
        }
        return render_template('tool.html', data=results)
@app.route('/search', methods=['GET', 'POST']) 
def search():
    app.secret_key = 'super secret key'
    search = SearchBar(request.form)
    if request.method == 'POST':
        return tool(search)
    return render_template('search.html', form=search)

@app.route('/')
def about():
    return render_template('index.html')

if __name__ == '__main__':
    app.run()