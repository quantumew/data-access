package daos

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/quantumew/data-access/models"
	"golang.org/x/net/context"
)

// JobDAO persists job data in database
type jobDAO struct{}

// NewJobDAO creates a new JobDAO
func NewJobDAO() *jobDAO {
	return &jobDAO{}
}

// Get reads the job with the specified ID from the database.
func (dao *jobDAO) Get(db *mongo.Database, name string) (*models.Job, error) {
	return dao.get(
		db,
		bson.NewDocument(
			bson.EC.String("name", name),
		),
	)
}

// GetByName reads the job with the specified name from the database.
func (dao *jobDAO) GetByName(db *mongo.Database, name string) (*models.Job, error) {
	return dao.get(
		db,
		bson.NewDocument(
			bson.EC.String("name", name),
		),
	)
}

func (dao *jobDAO) get(db *mongo.Database, doc *bson.Document) (*models.Job, error) {
	var job *models.Job
	col := db.Collection("job")
	result := bson.NewDocument()

	err := col.FindOne(
		context.Background(),
		doc,
	).Decode(result)

	if err != nil {
		return job, err
	}

	job, err = models.NewJobFromDoc(result)

	return job, err
}

// Create saves a new job record in the database.
// The Job.ID field will be populated with an automatically generated ID upon successful saving.
func (dao *jobDAO) Create(db *mongo.Database, job *models.Job) error {
	col := db.Collection("job")

	jobBson := models.NewDocFromJob(job)
	_, err := col.InsertOne(
		context.Background(),
		jobBson,
	)

	return err
}

// Update saves the changes to an job in the database.
func (dao *jobDAO) Update(db *mongo.Database, name string, job *models.Job) error {
	if _, err := dao.Get(db, name); err != nil {
		return err
	}

	jobBson := models.NewDocFromJob(job)
	col := db.Collection("job")
	_, err := col.UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("name", job.Name),
		),
		jobBson,
	)

	return err
}

// Delete deletes an job with the specified ID from the database.
func (dao *jobDAO) Delete(db *mongo.Database, name string) error {
	_, err := dao.Get(db, name)
	if err != nil {
		return err
	}

	col := db.Collection("job")
	_, err = col.DeleteOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("name", name),
		),
	)

	return err
}

// Count returns the number of the job records in the database.
func (dao *jobDAO) Count(db *mongo.Database) (int64, error) {
	col := db.Collection("job")

	return col.Count(
		context.Background(),
		bson.NewDocument(),
	)
}

// Query retrieves the job records with the specified offset and limit from the database.
func (dao *jobDAO) Query(db *mongo.Database, offset, limit int) ([]*models.Job, error) {
	jobList := []*models.Job{}
	col := db.Collection("job")
	ctx := context.Background()

	cursor, err := col.Find(
		ctx,
		bson.NewDocument(),
	)
	defer cursor.Close(ctx)
	elm := bson.NewDocument()

	for cursor.Next(ctx) {
		elm.Reset()

		if err := cursor.Decode(elm); err != nil {
			return jobList, err
		}
		job, err := models.NewJobFromDoc(elm)

		if err != nil {
			return jobList, err
		}

		jobList = append(jobList, job)
	}

	return jobList, err
}
