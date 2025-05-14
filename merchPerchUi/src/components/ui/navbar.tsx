// components/NavBar.tsx
import { Link, useLocation } from 'react-router-dom'

const NavBar: React.FC = () => {
  const location = useLocation()

  return (
    <nav className="flex gap-4 p-4 border-b">
      <Link
        to="/"
        className={location.pathname === '/' ? 'font-bold underline' : ''}
      >
        Artists
      </Link>
      <Link
        to="/products"
        className={location.pathname === '/products' ? 'font-bold underline' : ''}
      >
        Products
      </Link>
    </nav>
  )
}

export default NavBar
