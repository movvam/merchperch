
// import React from 'react'
import {
  Card,
  CardContent
} from "./ui/card"
import { cn } from "@/lib/utils"
import { WindowScroller, AutoSizer, Grid } from 'react-virtualized'

// maybe make these dynamic in the future
const COLUMN_WIDTH = 300
const ROW_HEIGHT = 320
const GAP = 16 // Optional spacing
const columnCount = 3 // Or calculate based on screen width dynamically

export type ProductData = {
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

type ProductCardProps = {
  productData: ProductData
}

const ProductCard: React.FC<ProductCardProps> = ({ productData }) => {
  return (
    <div className="h-60 pb-60 cursor-pointer bg-cover bg-center border rounded">
      <img
        src={productData.product.images.edges[0]?.node.url}
        alt={productData.product.title}
        className="w-full h-40 object-cover"
      />
      <div className="p-2">
        <p className="spotify-product-name-text">{productData.product.title}</p>
      </div>
    </div>
  )
}

type ProductDisplayProps = {
  products: ProductData[]
}

function ProductJsonDataDisplay(props: ProductDisplayProps) {
  const products = props.products
  return(
<WindowScroller>
  {({ height, isScrolling, onChildScroll, scrollTop }) => (
    <AutoSizer disableHeight>
      {({ width }) => {
        const columnCount = Math.floor(width / (COLUMN_WIDTH + GAP))
        const rowCount = Math.ceil(products.length / columnCount)

        return (
          <Grid
            autoHeight
            height={height}
            width={width}
            columnWidth={COLUMN_WIDTH}
            columnCount={columnCount}
            rowHeight={ROW_HEIGHT}
            rowCount={rowCount}
            isScrolling={isScrolling}
            onScroll={onChildScroll}
            scrollTop={scrollTop}
            cellRenderer={({ columnIndex, rowIndex, key, style }) => {
              const index = rowIndex * columnCount + columnIndex
              const product = products[index]
              if (!product) return null

              return (
                <div key={key} style={{ ...style, padding: GAP / 2 }}>
                  <ProductCard productData={product} />
                </div>
              )
            }}
          />
        )
      }}
    </AutoSizer>
  )}
</WindowScroller>
  )
}

export default ProductJsonDataDisplay;
