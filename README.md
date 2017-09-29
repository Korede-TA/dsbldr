# dsbldr
Golang tool that makes it easier to build Machine Learning Datasets from public Social Network APIs. It's fairly low level; primarily helps out with managing concurrent operations and saving data under the hood -- while providing a straightforward and intuitive API to do so.

## Abstract
(Don't know if people actually write these for OSS projects)

Social APIs are a common target for datasets for training Machine Learning models. This is exemplified by models such as [Tweet2Vec](https://arxiv.org/abs/1607.07514) that directly try to extract features or create embeddings from or for this kind of data. Other kind of models that benefit from this kind of data are NLP-oriented models which might just benefit from a large repo of structured (or not so structured) text that may or may not have some kind of labeling.

It seems that people often revert to specific publicly available general datasets to get such information and often have to settle for stuff that's not specific enough. Alternatively, one might revert to putting together an elaborate feature extraction pipeline which isn't really fun either and can take a lot of time out of actual feature engineering and model formulation work. I had faced a myraid of challenges along these veins when trying to put together a dataset from Twitter data and wondered why there wasn't already some kind of tool to do this.

I'm proposing a rather simple solution that, I think would be pretty useful for fairly simple datasets using a Feature-based API to retrieve specific features through requests to different endpoints on an API and writing this all to a desired data format.

To be totally honest, I'm very new to the Machine Learning space but I've been learning a lot very quickly and I'm very open to feedback. I've barely started actually writing this but I've been brewing it quite thoroughly in my head and I hope to be able to share something genuinely useful in the coming weeks :)

## TODO
- [x] Top level Feature-based API
- [ ] Concurrency stuff using Goroutines, channels and all that fun stuff
- [ ] Saving functionality for different data formats
    - [x] CSV
    - [ ] JSON
    - ...
- [ ] Support for different API data formats
    - [ ] JSON
    - [ ] XML
    - ...
- [ ] Authentication
- [ ] Command line functionality
- [ ] Demo
- (Other stuff as it comes to mind, please feel free to make suggestions) ...