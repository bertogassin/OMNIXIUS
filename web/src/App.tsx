import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import { AuthProvider } from './auth';
import Layout from './Layout';
import ProtectedRoute from './ProtectedRoute';
import Login from './pages/Login';
import Register from './pages/Register';
import ForgotPassword from './pages/ForgotPassword';
import Dashboard from './pages/Dashboard';
import Marketplace from './pages/Marketplace';
import Product from './pages/Product';
import ProductCreate from './pages/ProductCreate';
import Orders from './pages/Orders';
import Order from './pages/Order';
import FindProfessional from './pages/FindProfessional';
import Notifications from './pages/Notifications';
import Profile from './pages/Profile';
import Settings from './pages/Settings';
import Wallet from './pages/Wallet';
import Remittances from './pages/Remittances';
import Vault from './pages/Vault';
import Trade from './pages/Trade';
import Connect from './pages/Connect';
import AI from './pages/AI';
import Mail from './pages/Mail';
import Conversation from './pages/Conversation';
import Admin from './pages/Admin';
import Direction from './pages/Direction';
import ProfileEdit from './pages/ProfileEdit';
import ResetPassword from './pages/ResetPassword';

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter basename="/app">
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<ProtectedRoute><Dashboard /></ProtectedRoute>} />
            <Route path="marketplace" element={<ProtectedRoute><Marketplace /></ProtectedRoute>} />
            <Route path="product/:id" element={<ProtectedRoute><Product /></ProtectedRoute>} />
            <Route path="product-create" element={<ProtectedRoute><ProductCreate /></ProtectedRoute>} />
            <Route path="orders" element={<ProtectedRoute><Orders /></ProtectedRoute>} />
            <Route path="order/:id" element={<ProtectedRoute><Order /></ProtectedRoute>} />
            <Route path="find-professional" element={<ProtectedRoute><FindProfessional /></ProtectedRoute>} />
            <Route path="notifications" element={<ProtectedRoute><Notifications /></ProtectedRoute>} />
            <Route path="profile" element={<ProtectedRoute><Profile /></ProtectedRoute>} />
            <Route path="settings" element={<ProtectedRoute><Settings /></ProtectedRoute>} />
            <Route path="wallet" element={<ProtectedRoute><Wallet /></ProtectedRoute>} />
            <Route path="remittances" element={<ProtectedRoute><Remittances /></ProtectedRoute>} />
            <Route path="vault" element={<ProtectedRoute><Vault /></ProtectedRoute>} />
            <Route path="trade" element={<ProtectedRoute><Trade /></ProtectedRoute>} />
            <Route path="connect" element={<ProtectedRoute><Connect /></ProtectedRoute>} />
            <Route path="ai" element={<ProtectedRoute><AI /></ProtectedRoute>} />
            <Route path="mail" element={<ProtectedRoute><Mail /></ProtectedRoute>} />
            <Route path="conversation/:id" element={<ProtectedRoute><Conversation /></ProtectedRoute>} />
            <Route path="admin" element={<ProtectedRoute><Admin /></ProtectedRoute>} />
            <Route path="direction/:name" element={<ProtectedRoute><Direction /></ProtectedRoute>} />
            <Route path="profile-edit" element={<ProtectedRoute><ProfileEdit /></ProtectedRoute>} />
            <Route path="login" element={<Login />} />
            <Route path="register" element={<Register />} />
            <Route path="forgot-password" element={<ForgotPassword />} />
            <Route path="reset-password" element={<ResetPassword />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}
