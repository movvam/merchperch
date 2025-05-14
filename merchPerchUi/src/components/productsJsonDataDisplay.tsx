
// import React from 'react'
import JsonData from '../data/products.json'
import {
  Card,
  CardContent
} from "./ui/card"
import { cn } from "@/lib/utils"

type ProductData = {
  product: {
    id: string;
    title: string;
    productType: string;
    handle: string;
    images: {
      edges: {
        node: {
          url: string;
        };
      }[];
    };
    variants: {
      edges: {
        node: {
          price: {
            amount: string;
            currencyCode: string;
          };
        };
      }[];
    };
  };
  artist: string;
};

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

function ProductJsonDataDisplay(){
	const DisplayData=JsonData.map(
		(productData: ProductData)=>{
			return(
        <Container>
            <Card 
            backgroundImage={productData.product.images.edges[0]?.node.url} 
            className="auto-rows-auto  h-60 pb-60 cursor-pointer">
              {/* <CardContent className="auto-rows-auto	"> */}
                <p className="spotify-product-name-text">
                  {productData.product.title}
                </p>
              {/* </CardContent> */}

            {/* </Card> */}/</Card>
        </Container>
			)
		}
	)

	return(
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
			{DisplayData}
		</div>
	)
}

export default ProductJsonDataDisplay;
