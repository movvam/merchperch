import ProductJsonDataDisplay, { ProductData } from "@/components/productsJsonDataDisplay";
import ProductJsonData from '../data/products.json'
import ArtistJsonData from '../data/artistShops.json'
import { useMemo, useState } from "react";

const products = ProductJsonData as ProductData[]
const artists = ArtistJsonData as any[]

const artistMap: Record<string, any> = artists.reduce((map, artist) => {
  map[artist.id] = artist
  return map
}, {} as Record<string, any>)


type SortOption = 'mainstreamFirst' | 'undergroundFirst' | 'priceLowToHigh' | 'priceHighToLow'
const sortOptions: { label: string; value: SortOption }[] = [
  { label: 'Mainstream First', value: 'mainstreamFirst' },
  { label: 'Underground First', value: 'undergroundFirst' },
  { label: 'Price: Low to high', value: 'priceLowToHigh' },
  { label: 'Price: High to low', value: 'priceHighToLow' },
]


function ProductsPage() {
  const [sortBy, setSortBy] = useState<SortOption>('mainstreamFirst')

  const artistMap: Record<string, any> = useMemo(() => {
    return artists.reduce((map, artist) => {
      map[artist.id] = artist
      return map
    }, {} as Record<string, any>)
  }, [artists])

  const sortedProducts = useMemo(() => {
    return [...products].sort((a, b) => {
      const priceA = parseFloat(a.product.variants.edges[0]?.node.price.amount || '0')
      const priceB = parseFloat(b.product.variants.edges[0]?.node.price.amount || '0')
      const popularityA = artistMap[a.artist]?.popularity || 0
      const popularityB = artistMap[b.artist]?.popularity || 0

      switch (sortBy) {
        case 'priceLowToHigh':
          return priceA - priceB
        case 'priceHighToLow':
          return priceB - priceA
        case 'mainstreamFirst':
          return popularityB - popularityA
        case 'undergroundFirst':
          return popularityA - popularityB
      }
    })
  }, [products, sortBy, artistMap])
  return (
    <div className="flex app">
      {/* Left: Categories */}
      <aside className="w-1/5 p-4">
        {/* Category filters here */}
      </aside>

      {/* Center: Products */}
      <main className="w-3/5 p-4">
        <ProductJsonDataDisplay products={sortedProducts}/> 
      </main>

      {/* Right: Sort options */}
      <aside className="w-1/5 p-4">
      <h6 className="font-semibold mb-2">sort by</h6>
        <ul className="space-y-1 text-sm">
          {sortOptions.map(({ label, value }) => (
            <li key={label}>
              <button className={`hover:underline ${sortBy === value ? 'font-bold' : 'font-normal'}`} 
                onClick={() => setSortBy(value)}>
                {label}
              </button>
            </li>
          ))}
        </ul>
      </aside>
    </div>
    )
}

export default ProductsPage

