import React, { useState } from 'react';
import { useAuthContext } from '../contexts/authContext';
import { getErrorMessage } from '../utils/errorHandler';
import NetworkSVG from '../assets/network.svg';
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

const Login: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLogin, setIsLogin] = useState(true);
  const { login, signup } = useAuthContext();
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    
    try {
      if (isLogin) {
        await login(email, password);
      } else {
        await signup(email, password);
      }
    } catch (error) {
      const errorMessage = getErrorMessage(error);
      setError(errorMessage);
      console.error('Authentication error:', error);
    } finally {
      setLoading(false);
    }
  };

  const toggleMode = () => {
    setIsLogin(!isLogin);
    setError('');
  };

  return (
    <div className='flex flex-col md:flex-row items-center justify-center min-h-screen p-4 gap-6 md:gap-8'>
      <div className="flex justify-center items-center w-full md:w-1/2 lg:w-2/5 p-2 md:p-4 order-1">
        <img 
          src={NetworkSVG} 
          alt="Devices illustration" 
          className="w-full max-w-[200px] md:max-w-[300px] lg:max-w-[400px] xl:max-w-[500px] h-auto" 
        />
      </div>
      
      <Card className="w-full max-w-md md:w-1/2 lg:w-2/5 order-2 md:order-2">
        <CardHeader className="pb-4">
          <CardTitle className="text-2xl">
            {isLogin ? 'Login to your account' : 'Create an account'}
          </CardTitle>
          <CardDescription>
            {isLogin 
              ? 'Enter your email below to login to your account' 
              : 'Enter your email below to create your account'
            }
          </CardDescription>
        </CardHeader>
        
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4 pb-6">
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="email@example.com"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                disabled={loading}
                className="h-11"
              />
            </div>
            
            <div >
              <Label htmlFor="password">Password</Label>
              <Input 
                id="password" 
                type="password" 
                required 
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={loading}
                className="h-11"
              />
            </div>
            
            {error && (
              <div className="text-red-700 text-md text-center py-2">
                {error}
              </div>
            )}
          </CardContent>
          
          <CardFooter className="flex-col gap-4 pt-2">
            <Button 
              type="submit"
              variant="default"
              className="w-full h-11 cursor-pointer" 
              disabled={loading}
            >
              {loading ? 'Loading...' : (isLogin ? 'Login' : 'Sign Up')}
            </Button>
            
            <Button 
              type="button" 
              variant="link" 
              onClick={toggleMode}
              disabled={loading}
              className="h-9 cursor-pointer"
            >
              {isLogin ? "Don't have an account? Sign Up" : "Already have an account? Login"}
            </Button>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
};

export default Login;