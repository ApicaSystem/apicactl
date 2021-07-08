
which gox >& /dev/null
if [ $? != 0 ]
then
  echo "  gox command not found, run below to bring it"
  echo "    'go get github.com/mitchellh/gox'"
  echo "  exit!"
  exit 1
fi	


VER="2.0.4-oss"
CM=`git log -1 | grep commit | cut -d " " -f 2`
dd=`date`
LOC="logiqctl-release-$VER"
ALIST="logiqctl_darwin_amd64 logiqctl_freebsd_386 logiqctl_freebsd_amd64 logiqctl_freebsd_arm  \
logiqctl_linux_386 logiqctl_linux_amd64 logiqctl_linux_arm logiqctl_linux_mips logiqctl_linux_mips64 \
logiqctl_linux_mips64le logiqctl_linux_mipsle logiqctl_linux_s390x logiqctl_netbsd_386 \
logiqctl_netbsd_amd64 logiqctl_netbsd_arm logiqctl_openbsd_386 logiqctl_openbsd_amd64 \
logiqctl_windows_386.exe logiqctl_windows_amd64.exe"

#gox -osarch="$ALIST"
gox 

echo "ALIST=$ALIST"
rm -fr $LOC
mkdir $LOC

echo \
"logiqctl release form private build

    pkg date: $dd
     version: $VER
      commit: $CM
     os-arch: $ALIST
 packaged by: $USER

" > $LOC/aaa-readme.txt

for i in \
logiqctl_darwin_amd64 logiqctl_freebsd_386 logiqctl_freebsd_amd64 \
logiqctl_freebsd_arm logiqctl_linux_386 logiqctl_linux_amd64 \
logiqctl_linux_arm logiqctl_linux_mips logiqctl_linux_mips64 \
logiqctl_linux_mips64le logiqctl_linux_mipsle logiqctl_linux_s390x \
logiqctl_netbsd_386 logiqctl_netbsd_amd64 logiqctl_netbsd_arm \
logiqctl_openbsd_386 logiqctl_openbsd_amd64 logiqctl_windows_386.exe \
logiqctl_windows_amd64.exe
do
	mv $i $LOC
done

zip ${LOC}.zip ${LOC}/*
rm -fr ${LOC}

