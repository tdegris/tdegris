package lessons

type lesson struct {
	Text string
	Code string
}

var data = []struct {
	Title   string
	Content []lesson
}{
	{
		Title: "Chapter 1",
		Content: []lesson{
			{
				Text: "Lesson 1.1</br>Some other</br>text.",
				Code: "Some code for 1.1\nfloat32\npackage",
			},
			{
				Text: "Lesson 1.2",
				Code: "Some code for 1.2",
			},
		},
	},
	{
		Title: "Chapter 2",
		Content: []lesson{
			{
				Text: "Lesson 2.1",
				Code: "Some code for 2.1",
			},
			{
				Text: "Lesson 2.2",
				Code: "Some code for 2.2",
			},
		},
	},
}

type (
	Chapter struct {
		Title   string
		Content []*Lesson
	}

	Lesson struct {
		Chapter *Chapter

		Text string
		Code string
	}
)

func New() []Chapter {
	chapters := make([]Chapter, len(data))
	for i, chap := range data {
		chapters[i] = Chapter{
			Title: chap.Title,
		}
		lessons := make([]*Lesson, len(chap.Content))
		for lessonI, lesson := range chap.Content {
			lessons[lessonI] = &Lesson{
				Chapter: &chapters[i],
				Text:    lesson.Text,
				Code:    lesson.Code,
			}
		}
		chapters[i].Content = lessons
	}
	return chapters
}
