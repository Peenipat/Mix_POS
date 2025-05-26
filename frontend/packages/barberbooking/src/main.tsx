import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { BrowserRouter } from 'react-router-dom'
import { Provider } from 'react-redux'
import { store } from './store/index.ts'


createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter> {/* มีเพื่อให้ useNavigate, Routes, Link ใช้งานได้*/}
      <Provider store={store}> {/* component เข้าถึง Redux store ได้  */}
        <App />
      </Provider>
    </BrowserRouter>
  </StrictMode >,
)
