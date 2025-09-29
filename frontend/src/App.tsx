import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuthContext } from './contexts/authContext';
import Login from './pages/Login';
import Home from './pages/Home';
import Devices from './pages/Devices';
import { ThemeProvider } from './providers/ThemeProvider';
import Layout from './components/navbar/Layout';
import Notifications from './pages/Notifications';
import DeviceRegistrationPage from './pages/DeviceRegistrationPage';

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user } = useAuthContext();
  return user ? <Layout>{children}</Layout> : <Navigate to="/login" />;
};

const AppRoutes: React.FC = () => {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
       <Route 
          path="/" 
          element={
            <ProtectedRoute>
              <Home />
            </ProtectedRoute>
          } 
        />
         <Route 
          path="/devices" 
          element={
            <ProtectedRoute>
              <Devices />
            </ProtectedRoute>
          } 
        />
        <Route
          path='/notifications'
          element={
            <ProtectedRoute>
              <Notifications />
            </ProtectedRoute>
          }
        />
         <Route 
          path="/register" 
          element={
            <ProtectedRoute>
              <DeviceRegistrationPage />
            </ProtectedRoute>
          } 
        />   
    </Routes>
  );
};

const App: React.FC = () => {
  return (
    <AuthProvider>
      <ThemeProvider defaultTheme="light" storageKey="vite-ui-theme">
        <Router>
          <div className="App">
            <AppRoutes />
          </div>
        </Router>
      </ThemeProvider>
    </AuthProvider>
  );
};

export default App;
