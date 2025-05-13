
// import React from 'react'
import JsonData from '../../data/output-new.json'
import { ParallaxProvider,Parallax } from 'react-scroll-parallax'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./card"
import { cn } from "@/lib/utils"

function Container({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        "flex items-center pb-0 justify-center [&>div]:w-full",
        className
      )}
      {...props}
    />
  )
}

function ArtistJsonDataDisplay(){
	const DisplayData=JsonData.map(
		(info)=>{
			return(
        <Container>
          				{/* <div className=" h-60 pb-60"> */}
                  {/* <img className="bg-photo"
                    src={info.photo_url}
                    alt={`${info.name}.`}/> */}
          <a href={"https://"+info.shopify_url} target="_blank" rel="noopener noreferrer">
            <Card backgroundImage={info.photo_url} className="auto-rows-auto  h-60 pb-60 cursor-pointer">
              <CardContent className="auto-rows-auto	">



                <text className="spotify-artist-name-text">
                      {info.name}
                    </text>
              </CardContent>

              {/* <CardFooter className="flex justify-between">
              </CardFooter> */}
            </Card>
          </a>
				{/* </div> */}
        </Container>
			)
		}
	)

	return(
		<div>
			{DisplayData}
		</div>
	)
}

export default ArtistJsonDataDisplay;
