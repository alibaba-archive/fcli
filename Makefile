clean:
	rm -rf bundles && mkdir bundles

binary:
	bash hack/build.sh

package: clean binary

upload:
	bash hack/upload.sh