# -*- encoding:utf-8 -*-

import os
import  hashlib
from sys import argv

def  generate_file_suminfo(inputfile,outputfile):
  '''inputfile 为需要生成验证信息的文件 
      outputfile需要把验证信息写入的文件
  '''
  md5 = hashlib.md5(open(inputfile,'rb').read()).hexdigest()
  sha1 = hashlib.sha1(open(inputfile,'rb').read()).hexdigest()
  sha256 = hashlib.sha256(open(inputfile,'rb').read()).hexdigest()
  sha384 = hashlib.sha384(open(inputfile,'rb').read()).hexdigest()
  sha512 = hashlib.sha512(open(inputfile,'rb').read()).hexdigest()

  print "MD5: %s\n" % md5
  print "SHA1: %s\n" % sha1
  print "SHA256: %s\n" % sha256
  print "SHA384: %s\n" % sha384
  print "SHA512: %s\n" % sha512

  f = open(outputfile,"w")
  f.write("MD5: %s\n" % md5)
  f.write("SHA1: %s\n" % sha1)
  f.write("SHA256: %s\n" % sha256)
  f.write("SHA384: %s\n"  % sha384)
  f.write("SHA512: %s\n" % sha512)
  f.flush()
  f.close()

if __name__ == "__main__":
    if len(argv) == 3:
        generate_file_suminfo(argv[1],argv[2])
    else:
        print "you need two argument!\n"
        print "Usage: \n"
        print "           python2 checksum.py test.txt test.txt.hash"