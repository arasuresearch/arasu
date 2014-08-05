Arasu   :  A Lightning Fast Web Framework
=====
<h3><strong> Note: </strong>arasu development only work on dart enabled browsers like dartium or dart enabled chrome browser.</h3>

Arasu is a Next Generation Full Stack Web framework written on Go language & Dart language.  

Features
========
* lightning fast, because of golang and dartlang
* use RDBMS and BIGDATA for serverside
* use IndexedDB and Angular Dart for clientside
* use TDD default by golang and dartlang 
* use BDD with selenium and Spinach (this is in alpha)
* automatic build system.

Installation
============
* Install <a href="http://golang.org">Golang</a>
* Install <a href="http://dartlang.org">Dartlang</a>
* Install <a href="http://www.mysql.com">Mysql</a>
* Install <a href="http://hbase.apache.org">Hbase</a>
* Install Arasu Framework

open terminal

`go get github.com/arasuresearch/arasu` 

Creating a New Arasu Project
============================
<pre>
Part 1 : creating scaffold for relational Database Management System aka RDBMS  

~$ arasu new demo
~$ cd demo
~$ arasu dstore create
~$ arasu generate scaffold User name pass:string age:integer dob:timestamp sex:bool
~$ arasu dstore migrate  

Now start the server:

~$ arasu serve

now visit http://localhost:4000/ on <i><a href="https://www.dartlang.org/tools/dartium">Dartium</a> or dart enabled chrome</i> browser. 
then visit 

`http://localhost:4000/users.html`

There you are !!!

Part 2 : creating scaffold for BigData (hbase)
stop the server by pressing CTRL+C

start bigdata...

~$ start-hbase.sh
~$ hbase thrift start

leave this terminal to run thrift deamon



~$ arasu dstore create --dstore bigdata

this will result in failure

unfortunately hbase thrift API is not supporting to create database through API
so we have to create it manually . for that 
~$ hbase shell
 >  create_namespace 'demo_development'

close hbase shell 

~$ arasu generate scaffold Profile Contact{Name:String,Age:int,Sex:bool} --dstore bigdata
~$ arasu dstore migrate --dstore bigdata

Now start the server:

~$ arasu serve

now visit `http://localhost:4000/profiles.html` on dartium

<p>lets dive into Tutorial  <a href="http://arasuframework.org">Arasu Framework</a>.</p>

</pre>

License
============================
<p>Released under the <a href="http://www.opensource.org/licenses/MIT">MIT License</a>.</p>   

