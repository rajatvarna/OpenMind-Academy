import React, { useEffect, useState } from 'react';
import { View, Text, StyleSheet, FlatList } from 'react-native';
import { useAuth } from '../context/AuthContext';

// Placeholder lesson data
const LESSONS = {
  '1': [{ id: 'l1', title: 'Welcome to Python' }, { id: 'l2', title: 'Variables and Types' }],
  '2': [{ id: 'l1', title: 'What is HTML?' }, { id: 'l2', title: 'Your First Web Page' }],
  '3': [{ id: 'l1', title: 'Introduction to Happiness' }, { id: 'l2', title: 'Misconceptions About Happiness' }],
};

const CourseScreen = ({ route }) => {
  const { courseId, courseTitle } = route.params;
  const [lessons, setLessons] = useState([]);
  const { completedLessons } = useAuth();

  useEffect(() => {
    // In a real app, you would fetch lessons for the courseId from your API
    // For now, we'll use our placeholder data
    setLessons(LESSONS[courseId] || []);
  }, [courseId]);

  const renderLesson = ({ item }) => {
    const isCompleted = completedLessons.has(parseInt(item.id.replace('l', '')));
    return (
      <View style={[styles.lessonItem, isCompleted && styles.completedLesson]}>
        <Text style={isCompleted && styles.completedText}>
          {isCompleted ? '✔' : '○'} {item.title}
        </Text>
      </View>
    );
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>{courseTitle}</Text>
      <Text style={styles.description}>
        This is a placeholder description for the course. Here you would show more details about what the user will learn.
      </Text>

      <Text style={styles.lessonsHeader}>Lessons</Text>
      <FlatList
        data={lessons}
        keyExtractor={item => item.id}
        renderItem={renderLesson}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 8,
  },
  description: {
    fontSize: 16,
    color: '#555',
    marginBottom: 24,
  },
  lessonsHeader: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 12,
  },
  lessonItem: {
    backgroundColor: 'white',
    padding: 16,
    marginBottom: 8,
    borderRadius: 6,
  },
  completedLesson: {
    backgroundColor: '#e8f5e9',
    borderColor: '#a5d6a7',
    borderWidth: 1,
  },
  completedText: {
    color: '#388e3c',
    textDecorationLine: 'line-through',
  }
});

export default CourseScreen;
