import React from 'react';
import { View, Text, FlatList, TouchableOpacity, StyleSheet } from 'react-native';
import { useAuth } from '../context/AuthContext';

// Placeholder data - in a real app, this would be fetched from the Content Service API
const COURSES = [
  { id: '1', title: 'Introduction to Python', description: 'Learn the basics of Python programming.' },
  { id: '2', title: 'Web Development Fundamentals', description: 'Understand HTML, CSS, and JavaScript.' },
  { id: '3', title: 'The Science of Well-Being', description: 'A course on the science of happiness.' },
];

const HomeScreen = ({ navigation }) => {
  const { user, stats } = useAuth();

  const renderItem = ({ item }) => (
    <TouchableOpacity
      style={styles.card}
      onPress={() => navigation.navigate('Course', { courseId: item.id, courseTitle: item.title })}
    >
      <Text style={styles.cardTitle}>{item.title}</Text>
      <Text>{item.description}</Text>
    </TouchableOpacity>
  );

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.headerText}>Welcome, {user.name}</Text>
        <Text style={styles.scoreText}>⭐ {stats.score} Points</Text>
      </View>
      <FlatList
        data={COURSES}
        renderItem={renderItem}
        keyExtractor={item => item.id}
        contentContainerStyle={styles.listContainer}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  header: {
    padding: 20,
    paddingBottom: 10,
    backgroundColor: 'white',
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
  },
  headerText: {
    fontSize: 22,
    fontWeight: 'bold',
  },
  scoreText: {
    fontSize: 16,
    color: '#555',
    marginTop: 4,
  },
  listContainer: {
    padding: 16,
  },
  card: {
    backgroundColor: 'white',
    padding: 20,
    marginBottom: 16,
    borderRadius: 8,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  cardTitle: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 8,
  },
});

export default HomeScreen;
