
// import React from 'react'
import JsonData from '../../data/shopData.json'
import {
  Card,
  CardContent
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
          <a href={"https://"+info.shopify_url} target="_blank" rel="noopener noreferrer">
            <Card backgroundImage={info.photo_url} className="auto-rows-auto  h-60 pb-60 cursor-pointer">
              <CardContent className="auto-rows-auto	">



                <text className="spotify-artist-name-text">
                      {info.name}
                    </text>
              </CardContent>

            </Card>
          </a>
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
