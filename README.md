# Docker apps getting started

A simple go project that can be deployed by docker-app on Liara. Notice that now Liara provides [go-platform](https://liara.ir/landing/%D9%87%D8%A7%D8%B3%D8%AA-%DA%AF%D9%88%D9%84%D9%86%DA%AF-golang/) and you can use it for deploying golang projects
and this repo is archived.

## Deployment

- Create new [docker app](https://console.liara.ir/apps/create) & install [Liara CLI](https://docs.liara.ir/cli/install)
- Set ENVs in `.env.example` on the docker app
- Clone this repo by `git clone https://github.com/liara-cloud/docker-go-getting-started.git`
- Navigate to `docker-go-getting-started` dir
- Deploy app using command `liara deploy --platform docker --port 8080`

And ... that's it. you can enjoy your app!

