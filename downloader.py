#coding=utf-8
from bs4 import BeautifulSoup
import requests
import os
import shutil
import subprocess
from progressbar import ProgressBar
import threadpool 
import threading
import time

def read_all_string(filename):
    with open(filename,'r',encoding='utf8') as f:
        data = f.read()
        return data
def write_all_string(filename,data):
    with open(filename,'w',encoding='utf8') as f:
        f.write(data)

def generate_txt(book):
        chapters = info['chapters']
        content = ""
        for index in range(len(chapters)):
                cinfo = chapters[index]
                content+= "第 " + str(index+1)+" 章 " + cinfo["title"] + "\n\n"
                content+= cinfo["content"].replace("readx();","")+"\n\n"
        outfname = "./outputs/" + info["name"] + ".txt"
        write_all_string(outfname,content)


def generate_mobi(book):
    tpl_cover = read_all_string("./tpls/tpl_cover.html")
    tpl_book_toc =read_all_string("./tpls/tpl_book_toc.html")
    tpl_chapter = read_all_string("./tpls/tpl_chapter.html")
    tpl_content = read_all_string("./tpls/tpl_content.opf")
    tpl_style = read_all_string("./tpls/tpl_style.css")
    tpl_toc = read_all_string("./tpls/tpl_toc.ncx")
    path = "./tmp/" + book['name']

    if os.path.exists(path) :
        shutil.rmtree(path)
    
    os.makedirs(path)
    # 封面
    cover = tpl_cover.replace("___BOOK_NAME___",book['name'])
    cover = cover.replace("___BOOK_AUTHOR___",book['author'])
    write_all_string(path+"/cover.html",cover)

    # 章节
    chapters = info['chapters']
    toc_content = ""
    nax_toc_content = ""
    opf_toc = ""
    opf_spine = ""
    for index in range(len(chapters)):
        cinfo = chapters[index]
        chapter = tpl_chapter.replace("___CHAPTER_ID___","Chapter " + str(index))
        chapter = chapter.replace("___CHAPTER_NAME___",cinfo['title'])

        content = cinfo['content']
        content = content.replace("readx();","")
        #content = content.replace('\r','')
        content_lines = content.split('\r')
        content = ""
        for line in content_lines:
                content += "<p class=\"a\">    {0}</p>".format(line)


        chapter = chapter.replace("___CONTENT___",content)
        cpath = path+"/chapter"+str(index)+".html"
        write_all_string(cpath,chapter)

        toc_line = "<dt class=\"tocl2\"><a href=\"chapter{0}.html\">{1}</a></dt>\n".format(index,cinfo["title"])
        toc_content+=toc_line

        nax_toc_line = "<navPoint id=\"chapter{0}\" playOrder=\"{1}\">\n".format(index,index+1)
        nax_toc_content+=nax_toc_line
        nax_toc_line = "<navLabel><text>{0}</text></navLabel>\n".format(cinfo["title"])
        nax_toc_content+=nax_toc_line
        nax_toc_line = "<content src=\"chapter{0}.html\"/>\n</navPoint>\n".format(index)
        nax_toc_content+=nax_toc_line

        opf_toc += "<item id=\"chapter{0}\" href=\"chapter{0}.html\" media-type=\"application/xhtml+xml\"/>\n".format(index)
        opf_spine += "<itemref idref=\"chapter{0}\" linear=\"yes\"/>\n".format(index)




    # style
    write_all_string(path+"/style.css",tpl_style)

    # 目录
    book_toc = tpl_book_toc.replace("___CONTENT___",toc_content)
    write_all_string(path+"/book-toc.html",book_toc)

    nax_toc = tpl_toc.replace("___BOOK_ID___","11111")
    nax_toc = nax_toc.replace("___BOOK_NAME___",info['name'])
    nax_toc = nax_toc.replace("___NAV___",nax_toc_content)
    write_all_string(path+"/toc.ncx",nax_toc)

    #opf
    opf_content = tpl_content.replace("___MANIFEST___",opf_toc)
    opf_content = opf_content.replace("___SPINE___",opf_spine)
    opf_content = opf_content.replace("___BOOK_ID___","11111")
    opf_content = opf_content.replace("___BOOK_NAME___",info['name'])
    write_all_string(path+"/content.opf",opf_content)

    if os.path.exists("./outputs") == False :   
        os.mkdir("./outputs")

    # 生成
    outfname = info["name"] + ".mobi"
    subprocess.call( os.getcwd() + "/tools/kindlegen.exe " + path+"/content.opf -c1 -o " + outfname , shell=True)

    # copy
    shutil.copyfile(path+"/"+outfname,"./outputs/"+outfname)


# 获取章节列表
def get_bookinfo(bookid):
    url = "https://www.xbiquge6.com/" + bookid + "/"
    r = requests.get(url)
    r.encoding = r.apparent_encoding
    soup = BeautifulSoup(r.text,features="lxml")
    bookinfo = {}
    # 书名
    sbookName = soup.find("meta",property="og:novel:book_name")
    bookinfo["name"] = sbookName["content"]
    sauthor = soup.find("meta",property="og:novel:author")
    bookinfo["author"] = sauthor["content"]
    # 章节列表
    l = soup.find("div",id="list")
    dds = l.find_all("dd")
    chapters = []
    for dd in dds:
        title = dd.a.get_text()
        path = dd.a["href"]
        chapterinfo = {"title":title,"link":path}
        chapters.append(chapterinfo)
    bookinfo["chapters"] = chapters
    
    return bookinfo

def get_chapter_content(link):
    url = "https://www.xbiquge6.com/" + link
    r = requests.get(url)
    r.encoding = r.apparent_encoding
    soup = BeautifulSoup(r.text,features="lxml")
    l = soup.find("div",id="content")
    r = ""
    for text in l.stripped_strings:
            r+="\t"+text+"\n\n"

    return r


g_lock = threading.Lock()
g_totals = 0
g_results = {}
g_pbar = ProgressBar()
def get_chapter_async(link):
        content = get_chapter_content(link)
        g_lock.acquire()
        g_results[link] = content
        g_pbar.update(int(((len(g_results)) / (g_totals - 1)) * 100))
        g_lock.release()

def download_chapters(info):
    print("正在下载")
    pool = threadpool.ThreadPool(50)
    chapters = info['chapters']  
    g_pbar.start()
    links = []
    for i in range(len(chapters)):
        links.append(chapters[i]['link'])

    requests = threadpool.makeRequests(get_chapter_async,links)
    [pool.putRequest(req) for req in requests]
    pool.wait()
    for i in range(len(chapters)):
        chapters[i]['content'] = g_results[chapters[i]['link']]


print("正在获取书籍信息...")
info = get_bookinfo("0_642")
print(info["name"])
g_totals = len(info["chapters"])
print("共有 " + str(len(info["chapters"])) + " 章")
download_chapters(info)
generate_txt(info)
