package main

/*
FileDropper
https://www.filedropper.com/
1  #!/usr/bin/env python
2  # -*- coding: utf-8 -*-
3  #
4  # Copyright (c) 2009, Thomas Jost <thomas.jost@gmail.com>
5  #
6  # Permission to use, copy, modify, and/or distribute this software for any
7  # purpose with or without fee is hereby granted, provided that the above
8  # copyright notice and this permission notice appear in all copies.
9  #
10  # THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
11  # WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
12  # MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
13  # ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
14  # WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
15  # ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
16  # OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
17
18  """Interact with FileDropper.com"""
19
20  import os.path, re, urllib, urllib2
21  from types import StringTypes
22
23  from BeautifulSoup import BeautifulSoup
24  import poster.encode, poster.streaminghttp
25
26  # Permissions
27  FD_PERM_PUBLIC   = 0
28  FD_PERM_PASSWORD = 1
29  FD_PERM_PRIVATE  = 2
30
31  # Prefix for all the URLs used in the module
32  FD_URL         = "http://www.filedropper.com/"
33  FD_LOGIN_URL   = FD_URL + "login.php"
34  FD_PREMIUM_URL = FD_URL + "premium.php"
35  FD_UPLOAD_URL  = FD_URL + "index.php?xml=true"
36  FD_PERM_URL    = FD_PREMIUM_URL + "?action=setpermissions&%(varset)s&id=%(id)d"
37  FD_DELETE_URL  = FD_PREMIUM_URL + "?action=delete&id=%d"
38
39  # Max allowed size
40  FD_MAX_SIZE = 5*(1024**3)
41
42  # Error codes
43  FD_ERROR_NOT_LOGGED_IN    = 1
44  FD_ERROR_PERMISSIONS      = 2
45  FD_ERROR_FILE_TOO_BIG     = 3
46  FD_ERROR_UPLOAD           = 4
47  FD_ERROR_INVALID_PASSWORD = 5
48  FD_ERROR_REQUEST          = 6
49
50  FD_ERRORS = {
51      FD_ERROR_NOT_LOGGED_IN   : "Not logged in",
52      FD_ERROR_PERMISSIONS     : "Invalid permissions",
53      FD_ERROR_FILE_TOO_BIG    : "File is too big",
54      FD_ERROR_UPLOAD          : "Error while uploading a file",
55      FD_ERROR_INVALID_PASSWORD: "Invalid password",
56      FD_ERROR_REQUEST         : "Problem with the request",
57  }
58
59  # Regexp used for parsing HTML pages
60  FD_RE_LIST = re.compile("^unhide\('(\d+)'\)")
61
62 -class FileDropperException(Exception):
63 -    def __init__(self, errno):
64          self.errno = errno
65          self.errmsg = FD_ERRORS[errno]
66
67 -    def __str__(self):
68          return "[FileDropper error %d] %s" % (self.errno, self.errmsg)
69
70 -class FileDropper:
71      """Builds an empty FileDropper object that may be used for uploading files
72      to FileDropper.com, get details about files in a premium account, change
73      their permission or delete them."""
74
75 -    def __init__(self):
76          self.logged_in = False
77
78          # Init the streaming HTTP handler
79          #register_openers()
80
81          # Init the URL opener
82          #self.url = urllib2.build_opener(urllib2.HTTPCookieProcessor())
83          self.url = urllib2.build_opener(
84                  poster.streaminghttp.StreamingHTTPHandler,
85                  poster.streaminghttp.StreamingHTTPRedirectHandler,
86                  urllib2.HTTPCookieProcessor
87          )
88
89 -    def __del__(self):
90          if self.logged_in:
91              self.logout()
92
93 -    def login(self, username, password):
94          """Log into the premium account with the given username and password."""
95
96          # Build the data string to send with the POST request
97          data = urllib.urlencode({"username": username, "password": password})
98
99          # Send the request
100          res = self.url.open(FD_LOGIN_URL, data)
101
102          # What is our final URL?
103          dst_url = res.geturl()
104
105          self.logged_in = (res.getcode() == 200) and (dst_url == FD_PREMIUM_URL)
106          return self.logged_in
107
108 -    def logout(self):
109          """Log out from a premium account."""
110
111          if not self.logged_in:
112              raise FileDropperException(FD_ERROR_NOT_LOGGED_IN)
113
114          self.url.open(FD_URL + "login.php?action=logout")
115
116 -    def list(self):
117          """Get a list of files in the file manager of a premium account.
118
119          The return value is a list of 7-value tuples of the form
120          (file_name, id, downloads, size, date, permissions, public_url)"""
121
122          if not self.logged_in:
123              raise FileDropperException(FD_ERROR_NOT_LOGGED_IN)
124
125          # Download the page
126          html = self.url.open(FD_PREMIUM_URL).read()
127
128          # Parse it
129          soup = BeautifulSoup(html)
130
131          # Get files info
132          tags = [tag.parent for tag in soup.findAll("a", onclick = FD_RE_LIST)]
133          files = []
134          for tag in tags:
135              # File name in the link
136              file_name = tag.a.string.strip()
137
138              # File ID found using the regexp
139              m = FD_RE_LIST.search(tag.a['onclick'])
140              file_id = int(m.group(1))
141
142              div = tag.div
143
144              # Some dirty searches in an ugly div section
145              downloads = int(div.contents[2].replace('|', '').strip())
146              size = div.contents[4].replace('|', '').strip()      #TODO: parse it correctly
147              date = div.contents[6].replace('&nbsp;', '').strip() #TODO: parse it correctly
148
149              # Permissions: conversion from 2 strings to a symbol
150              raw_perm = div.find("span", id="fileperms[%d]" % file_id)
151              permissions = -1
152              if raw_perm.span.string == "Private":
153                  permissions = FD_PERM_PRIVATE
154              elif raw_perm.span.string == "Public":
155                  # If there is a <b> tag, it contains "No password"
156                  if raw_perm.b is not None:
157                      permissions = FD_PERM_PUBLIC
158                  else:
159                      permissions = FD_PERM_PASSWORD
160              else:
161                  raise FileDropperException(FD_ERROR_PERMISSIONS)
162
163              # Public URL (may be published safely anywhere)
164              public_url = div.find("input", type="text")['value']
165
166              value = (file_name, file_id, downloads, size, date, permissions, public_url)
167              files.append(value)
168
169          return files
170
171 -    def upload(self, filename):
172          """Upload the specified file"""
173
174          # Check the file size
175          if os.path.getsize(filename) > FD_MAX_SIZE:
176              raise FileDropperException(FD_ERROR_FILE_TOO_BIG)
177
178          # Prepare the encoded data
179          base_name = os.path.basename(filename)
180
181          mp1 = poster.encode.MultipartParam("Filename", base_name)
182          mp2 = poster.encode.MultipartParam("file", filename=base_name, filetype="application/octet-stream", fileobj=open(filename))
183
184          data, headers = poster.encode.multipart_encode([mp1, mp2])
185
186          # Prepare the request
187          req = urllib2.Request(FD_UPLOAD_URL, data, headers)
188
189          # Send the request
190          res = self.url.open(req)
191
192          # Get the intermediate url
193          tmp_url = res.read()
194          #TODO: check if upload failed...
195
196          # Get the real file URL... and end with a 404 error :)
197          try:
198              res = self.url.open(FD_URL + tmp_url[1:])
199          except urllib2.HTTPError, exc:
200              if exc.code == 404:
201                  return exc.geturl()
202              else:
203                  raise exc
204
205          # We should not reach this point as there is supposed to be a 404 error
206          raise FileDropperException(FD_ERROR_UPLOAD)
207
208 -    def set_perm(self, file_id, perm, password=None):
209          """Set new permissions for the specified file"""
210
211          if not self.logged_in:
212              raise FileDropperException(FD_ERROR_NOT_LOGGED_IN)
213
214          # Prepare the query
215          query = {'id': file_id}
216
217          # Set to public
218          if perm == FD_PERM_PUBLIC:
219              query['varset'] = 'public=true'
220
221              # There's a weird bug when changing from password-protected
222              # to public: permissions don't get updated unless we change
223              # them to private first
224              self.set_perm(file_id, FD_PERM_PRIVATE)
225
226          # Set to private
227          elif perm == FD_PERM_PRIVATE:
228              query['varset'] = 'private=true'
229
230          # Set to password-protected
231          elif perm == FD_PERM_PASSWORD:
232              if (type(password) not in StringTypes) or (password.strip() == ""):
233                  raise FileDropperException(FD_ERROR_INVALID_PASSWORD)
234              query['varset'] = urllib.urlencode({'password': password})
235
236          # Invalid case
237          else:
238              raise FileDropperException(FD_ERROR_PERMISSIONS)
239
240          # Prepare the request
241          url = FD_PERM_URL % query
242          res = self.url.open(url)
243
244          txt = res.read()
245          if res.getcode() != 200:
246              raise FileDropperException(FD_ERROR_REQUEST)
247
248          return txt
249
250 -    def delete(self, file_id):
251          """Delete the specified file"""
252
253          if not self.logged_in:
254              raise FileDropperException(FD_ERROR_NOT_LOGGED_IN)
255
256          # Do the query
257          res = self.url.open(FD_DELETE_URL % file_id)
258
259          txt = res.read()
260          if res.getcode() != 200:
261              raise FileDropperException(FD_ERROR_REQUEST)
262
263          return txt
264
265
266  if __name__ == "__main__":
267      from getpass import getpass
268      import sys
269
270      fd = FileDropper()
271      user = raw_input("Username: ")
272      if user != "":
273          password = getpass()
274          if not fd.login(user, password):
275              print "Login failed"
276
277      print "Current files:"
278      print fd.list()
279      print
280
281      uploaded_file = fd.upload("test.txt")
282      print "Upload: %s" % uploaded_file
283      print
284
285      print "New files:"
286      lst = fd.list()
287      print lst
288      file_id = -1
289      for file_data in lst:
290          if file_data[6] == uploaded_file:
291              file_id = file_data[1]
292              break
293      if file_id == -1:
294          print "Can't find file ID :-("
295          sys.exit(1)
296      print
297
298      print "Making the file private:"
299      fd.set_perm(file_id, FD_PERM_PRIVATE)
300      print fd.list()
301      print
302
303      print "Making the file password-protected:"
304      fd.set_perm(file_id, FD_PERM_PASSWORD, "passtest")
305      print fd.list()
306      print
307
308      print "Making the file public again:"
309      fd.set_perm(file_id, FD_PERM_PUBLIC)
310      print fd.list()
311      print
312
313      print "Making the file private again:"
314      fd.set_perm(file_id, FD_PERM_PRIVATE)
315      print fd.list()
316      print
317
318      print "Deleting the file:"
319      fd.delete(file_id)
320      print fd.list()
321    */
