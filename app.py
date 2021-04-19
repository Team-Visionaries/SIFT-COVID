import json
import newspaper #https://newspaper.readthedocs.io/en/latest/
from newspaper import Article 
from search import SearchBar
from bs4 import BeautifulSoup # https://www.crummy.com/software/BeautifulSoup/bs4/doc/
from flask import Flask, flash, render_template, request, redirect
from urllib.parse import urlparse
from datetime import datetime
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
        source_url = 'https://' + urlparse(url).netloc
        source = newspaper.build(source_url, memoize_articles=False)
        publish_date = article.publish_date.strftime("%d %B, %Y")
        diff = str((datetime.now() - article.publish_date).days)
        keyword_query = '+'.join(map(str, article.keywords))
        results = {
            'url': url,
            'title': article.title,
            'date': publish_date,
            'source_url': source_url,
            'source_name': source.brand.upper(),
            'source_descr': source.description,
            'keywords': ', '.join(map(str, article.keywords)),
            'summary': article.summary,
            'authors': ', '.join(map(str, article.authors)),
            'date_diff': diff,
            'google_query': 'https://www.google.com/search?q=' + keyword_query,
            'duck_query': 'https://www.duckduckgo.com?q=' + keyword_query,
            'yahoo_query': 'https://www.search.yahoo.com/search?p=' + keyword_query
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