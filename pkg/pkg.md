All latest packages should be stored here. 

To package the application, run the following script: 
```bash
cd pkg
fyne package \
  --os darwin \
  --icon ../img/icon.png \
  --name "Anba Abraam Service" \
  --app-version "1.0.0" \
  --release \
  --src ../cmd
```