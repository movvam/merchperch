
// import React from 'react'
import JsonData from '../../data/output.json'
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
          <Card backgroundImage={info.photo_url} className="auto-rows-auto  h-60 pb-60">
            <CardContent className="auto-rows-auto	">



              <text className="spotify-artist-name-text">
                    {info.name}
                  </text>
            </CardContent>

            {/* <CardFooter className="flex justify-between">
            </CardFooter> */}
          </Card>
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
