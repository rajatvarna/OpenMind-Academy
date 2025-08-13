import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { AuthProvider } from './context/AuthContext';

// Import screen components (we will create these next)
import HomeScreen from './screens/HomeScreen';
import CourseScreen from './screens/CourseScreen';
import LoginScreen from './screens/LoginScreen';

const Stack = createNativeStackNavigator();

function App() {
  return (
    <AuthProvider>
      <NavigationContainer>
        <Stack.Navigator initialRouteName="Home">
          <Stack.Screen
          name="Home"
          component={HomeScreen}
          options={{ title: 'Courses' }}
        />
        <Stack.Screen
          name="Course"
          component={CourseScreen}
          // The title can be set dynamically based on route params
          options={({ route }) => ({ title: route.params.courseTitle || 'Course Details' })}
        />
        <Stack.Screen
          name="Login"
          component={LoginScreen}
        />
      </Stack.Navigator>
    </NavigationContainer>
    </AuthProvider>
  );
}

export default App;
