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
def tool(url):
    if url == '' or not url.startswith('http'):
            flash('Please enter a URL!')
            return redirect('/search')
    else:
        errorMsg = "Could not find. Please open the article to search for this value"
        try:
            # Initiate article parsing
            article = Article(url)
            article.download()
            article.parse()
            article.nlp()
        except Exception:
            flash('There was an error parsing the article. Please try another article!')
            return redirect('/search')
        
        source_url = 'https://' + urlparse(url).netloc

        # Keywords are used to generate related search links
        # If keywords aren't available then the article title is used
        if len(article.keywords) != 0:
            search_query = '+'.join(map(str, article.keywords))
        else:
            search_query = article.title.replace(' ', '+')

        # Extracting the source from the URL
        try:
            source = newspaper.build(source_url, memoize_articles=False)
            source_brand = source.brand.upper()
        except Exception:
            source_brand = errorMsg

        # Extracting date the article was published
        try:
            publish_date = article.publish_date.strftime("%d %B, %Y")
            diff = str((datetime.now() - article.publish_date).days)
        except Exception:
            publish_date = errorMsg
            diff = "unknown"

        # Extracting citations citations
        citations = getCitations(article)
        if type(citations) == str:
            citations_len = str(0)
            citations = []
        else:
            citations_len = str(len(citations))

        results = {
            'url': url,
            'title': article.title,
            'date': publish_date,
            'source_url': source_url,
            'source_name': source_brand,
            'source_descr': source.description,
            'keywords': ', '.join(map(str, article.keywords)),
            'summary': article.summary,
            'authors': ', '.join(map(str, article.authors)),
            'date_diff': diff,
            'google_query': googleLink(search_query),
            'duck_query': duckLink(search_query),
            'yahoo_query': yahooLink(search_query),
            'bing_query': bingLink(search_query),
            'google_source': googleLink(source_brand + '+wikipedia'),
            'yahoo_source': yahooLink(source_brand + '+wikipedia'),
            'duck_source': duckLink(source_brand + '+wikipedia'),
            'bing_source': bingLink(source_brand + '+wikipedia'),
            'num_citations': citations_len,
            'citations': citations
        }
        return render_template('tool.html', data=results)
@app.route('/search', methods=['GET', 'POST']) 
def search():
    app.secret_key = 'super secret key'
    search = SearchBar(request.form)
    if request.args.get('url') != None:
        return tool(request.args.get('url'))
    if request.method == 'POST':
        return tool(search.data['search'])
    return render_template('search.html', form=search)

@app.route('/')
def about():
    return render_template('index.html')

def getCitations(article):
    try:
        htmlCode = article.html
        soup = BeautifulSoup(htmlCode, "html.parser")
        a_elems = soup.body.find_all("a")
        links = []
        for elem in a_elems:
            if(len(elem.attrs.keys()) <= 1 & ('href' in elem.attrs)):
                if elem.attrs['href'].startswith('http'):
                    links.append({'href': elem.attrs['href'], 'text': elem.text})
        return links
    except Exception:
        return "Could not find citations"

def googleLink(query):
    return 'https://www.google.com/search?q=' + query

def duckLink(query):
    return 'https://www.duckduckgo.com?q=' + query

def yahooLink(query):
    return 'https://www.search.yahoo.com/search?p=' + query

def bingLink(query):
    return 'https://www.bing.com/search?q=' + query
    
if __name__ == '__main__':
    app.run()