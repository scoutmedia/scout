# Scout
Scout is a media task scraping / fetching  service that is paired with frontend that allows users to search for movies / tv shows and select what media they want to add to their 
plex server. Each service is containerized and orchestrated via docker. 

## Inspiration
This was a hobby project mainly to scratch the itch and learn more about plex and creating my own server. During the process I found an old pc of mine , slapped linux on it and began 
dabbling with building a server. During this process I wondered what it would be like to just automate the process of adding new media to my plex server and that's how we got here.

## Goal
The goal of the project is create a service that automates the process of retrieving new media , this is accomplished by having our 
frontend service listening for users request from our  built using NextJS. The request is sent to the our scraper service
which scours the web for the requested media , fetches it and moves it to the respective file. 


## The Result

![terminal](https://media.discordapp.net/attachments/794650924115034133/1182058654334210189/Screenshot_from_2023-12-06_15-31-52.png?ex=658350e8&is=6570dbe8&hm=722326539b289be3ed7a6d1df1ed9419a6d2ff039b71b0fcff0996fa465a39a0&=&format=webp&quality=lossless)

Serving up Scout API to listen to incoming request

![mobile](https://media.discordapp.net/attachments/794650924115034133/1182058675909701632/Screenshot_from_2023-12-06_14-53-44.png?ex=658350ed&is=6570dbed&hm=b6ff5a3cb006d6c7d2de451b6854b651d6df776546a44f3f2ad3e0c100f08049&=&format=webp&quality=lossless&width=333&height=541)
![mobile](https://cdn.discordapp.com/attachments/794650924115034133/1182312395813703720/IMG_1275.png?ex=65843d39&is=6571c839&hm=37c6b3376abe756d687a5cb9e0c0ce40fc03b24fdf5726d0c9d4b7d1ed5f78aa&)
Mobile interface view for searching and requesting media.
