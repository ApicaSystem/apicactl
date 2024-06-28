# Release Checklist

## Steps

1. **Update Version in Code**
   - Update the version number in `cmd/root.go`.

2. **Update Documentation**
   - Update the "Release Notes" section in `README.md`.

3. **Update Build Script**
   - Update the version in `goxbuilt.sh`.

4. **Commit Changes**
   - Check into the `master` branch.

5. **Tag the Release**
   - Tag the release with the following commands:
     ```sh
     git tag -a <version> -m "released <version>, tagged on <date-time>"
     git push origin --tags
     ```

6. **Release on GitHub**
   - Use the GitHub UI to publish the release.

7. **Build Binaries**
   - Run the build script:
     ```sh
     sh goxbuilt.sh
     ```

8. **Upload Multi-Arch Binaries**
   - Upload the multi-architecture binaries.
