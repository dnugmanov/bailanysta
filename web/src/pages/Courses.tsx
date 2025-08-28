import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { BookOpen, Clock, Users, Star } from 'lucide-react'
import { apiClient } from '@/api/client'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

export default function CoursesPage() {
  const [activeTab, setActiveTab] = useState('all')

  const { data: coursesData, isLoading } = useQuery({
    queryKey: ['courses'],
    queryFn: () => apiClient.getCourses(),
  })

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold">Курсы</h1>
        </div>
        
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <Card key={i} className="animate-pulse">
              <CardHeader>
                <div className="h-4 bg-muted rounded w-3/4"></div>
                <div className="h-3 bg-muted rounded w-1/2"></div>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div className="h-3 bg-muted rounded"></div>
                  <div className="h-3 bg-muted rounded w-5/6"></div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    )
  }

  const courses = coursesData?.courses || []

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-6"
    >
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Курсы</h1>
        <Button onClick={() => {
          // Пока просто уведомление, потом можно сделать модальное окно
          alert('Функция создания курсов будет добавлена позже!')
        }}>
          <BookOpen className="h-4 w-4 mr-2" />
          Создать курс
        </Button>
      </div>

      <div className="flex space-x-2">
        {['all', 'my', 'popular'].map((tab) => (
          <Button
            key={tab}
            variant={activeTab === tab ? 'default' : 'outline'}
            onClick={() => setActiveTab(tab)}
          >
            {tab === 'all' && 'Все курсы'}
            {tab === 'my' && 'Мои курсы'}
            {tab === 'popular' && 'Популярные'}
          </Button>
        ))}
      </div>

      {courses.length === 0 ? (
        <Card className="p-12 text-center">
          <CardContent>
            <BookOpen className="h-12 w-12 mx-auto mb-4 text-muted-foreground" />
            <h3 className="text-lg font-semibold mb-2">Пока курсов нет</h3>
            <p className="text-muted-foreground mb-4">
              Станьте первым, кто создаст образовательный курс
            </p>
            <Button onClick={() => {
              alert('Функция создания курсов будет добавлена позже!')
            }}>
              <BookOpen className="h-4 w-4 mr-2" />
              Создать первый курс
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {courses.map((course) => (
            <motion.div
              key={course.id}
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              whileHover={{ scale: 1.02 }}
              transition={{ duration: 0.2 }}
            >
              <Card className="h-full hover:shadow-lg transition-shadow">
                <CardHeader>
                  <CardTitle className="line-clamp-2">{course.title}</CardTitle>
                  <p className="text-sm text-muted-foreground line-clamp-2">
                    {course.description}
                  </p>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between text-sm text-muted-foreground mb-4">
                    <div className="flex items-center">
                      <Clock className="h-4 w-4 mr-1" />
                      12 уроков
                    </div>
                    <div className="flex items-center">
                      <Users className="h-4 w-4 mr-1" />
                      24 студента
                    </div>
                    <div className="flex items-center">
                      <Star className="h-4 w-4 mr-1 fill-current text-yellow-500" />
                      4.8
                    </div>
                  </div>
                  
                  <Button 
                    className="w-full"
                    onClick={() => {
                      // Пока просто ссылка на курс, потом можно сделать отдельную страницу
                      window.open(`/course/${course.id}`, '_blank')
                    }}
                  >
                    Начать изучение
                  </Button>
                </CardContent>
              </Card>
            </motion.div>
          ))}
        </div>
      )}
    </motion.div>
  )
}
