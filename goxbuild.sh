
which gox >& /dev/null
if [ $? != 0 ]
then
  echo "  gox command not found, run below to bring it"
  echo "    'go get github.com/mitchellh/gox'"
  echo "  exit!"
  exit 1
fi	


VER=`grep currentReleaseVersion cmd/root.go | cut -d " " -f 4 | cut -d '"' -f 2`
echo "VER=$VER"

CM=`git log -1 | grep commit | cut -d " " -f 2`
dd=`date`
LOC="apicactl-release-$VER"
ALIST="apicactl_darwin_amd64 apicactl_freebsd_386 apicactl_freebsd_amd64 apicactl_freebsd_arm  \
apicactl_linux_386 apicactl_linux_amd64 apicactl_linux_arm apicactl_linux_mips apicactl_linux_mips64 \
apicactl_linux_mips64le apicactl_linux_mipsle apicactl_linux_s390x apicactl_netbsd_386 \
apicactl_netbsd_amd64 apicactl_netbsd_arm apicactl_openbsd_386 apicactl_openbsd_amd64 \
apicactl_windows_386.exe apicactl_windows_amd64.exe"

#gox -osarch="$ALIST"
gox 

echo "ALIST=$ALIST"
rm -fr $LOC
mkdir $LOC

echo \
"apicactl release form private build

    pkg date: $dd
     version: $VER
      commit: $CM
     os-arch: $ALIST
 packaged by: $USER

" > $LOC/aaa-readme.txt

for i in \
apicactl_darwin_amd64 apicactl_freebsd_386 apicactl_freebsd_amd64 \
apicactl_freebsd_arm apicactl_linux_386 apicactl_linux_amd64 \
apicactl_linux_arm apicactl_linux_mips apicactl_linux_mips64 \
apicactl_linux_mips64le apicactl_linux_mipsle apicactl_linux_s390x \
apicactl_netbsd_386 apicactl_netbsd_amd64 apicactl_netbsd_arm \
apicactl_openbsd_386 apicactl_openbsd_amd64 apicactl_windows_386.exe \
apicactl_windows_amd64.exe
do
	mv $i $LOC
done

zip ${LOC}.zip ${LOC}/*
rm -fr ${LOC}

