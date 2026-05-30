package repository

import (
	"context"
	"fmt"
	"time"

	"api-telemetria-robo/dto"
	"api-telemetria-robo/entity"
	"api-telemetria-robo/logs"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var pkgName logs.PackageName = "REPOSITORY"

type readingBSON struct {
	Timestamp string `bson:"timestamp"`
	Value     int    `bson:"value"`
}

type sensorBSON struct {
	Name     string        `bson:"name"`
	Readings []readingBSON `bson:"readings"`
}

type roundBSON struct {
	RoundNumber int          `bson:"roundNumber"`
	Sensors     []sensorBSON `bson:"sensors"`
}

type matchBSON struct {
	ID           int    `bson:"_id"`
	Title        string `bson:"title"`
	Date         string `bson:"date"`
	OpponentName string `bson:"opponentName"`
	Closed       bool   `bson:"isClosed"`
}

type MatchMongo struct {
	client     *mongo.Client
	db         string
	collection string
}

func NewMatchMongo(client *mongo.Client, db, collection string) *MatchMongo {
	return &MatchMongo{
		client:     client,
		db:         db,
		collection: collection,
	}
}

func (m *MatchMongo) CreateMatch(title string, date time.Time, opponentName string) error {
	var (
		collection = m.client.Database(m.db).Collection(m.collection)
		newMatch   matchBSON
		err        error
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := m.getNextIDCount(ctx)
	if err != nil {
		return err
	}

	newMatch = matchBSON{
		ID:           id,
		Title:        title,
		Date:         date.Format("2006-01-02"),
		OpponentName: opponentName,
		Closed:       false,
	}

	result, err := collection.InsertOne(ctx, newMatch)
	if err != nil {
		return err
	}

	logs.Infof(pkgName, "Created new match with ID: %v", result.InsertedID)
	return nil
}

func (m *MatchMongo) FindMatchByID(matchID int) (*dto.MatchDTO, error) {
	var (
		collection = m.client.Database(m.db).Collection(m.collection)
		result     matchBSON
		match      *dto.MatchDTO
		err        error
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": matchID}
	if err = collection.FindOne(ctx, filter).Decode(&result); err != nil {
		return nil, err
	}

	match, err = m.convertMatchBsonToDto(result)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func (m *MatchMongo) GetOpenMatch() (*dto.MatchDTO, error) {
	var (
		collection = m.client.Database(m.db).Collection(m.collection)
		result     matchBSON
		match      *dto.MatchDTO
		err        error
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	latestMatchID, err := m.getCurrentIDCount(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":      latestMatchID,
		"isClosed": false,
	}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	match, err = m.convertMatchBsonToDto(result)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func (m *MatchMongo) CloseMatch(matchID int) error {
	var (
		collection = m.client.Database(m.db).Collection(m.collection)
		err        error
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": matchID}
	update := bson.M{
		"$set": bson.M{
			"isClosed": true,
		},
	}
	if _, err = collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	if err = m.resetRoundCount(ctx); err != nil {
		return err
	}

	logs.Infof(pkgName, "Closed match with ID %d", matchID)
	return nil
}

func (m *MatchMongo) CreateNewRound(sensors []*dto.SensorDTO) error {
	var (
		collection  = m.client.Database(m.db).Collection(m.collection)
		openMatch   *dto.MatchDTO
		sensorsBSON []sensorBSON
		roundNumber int
		err         error
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if openMatch, err = m.GetOpenMatch(); err != nil {
		return err
	}

	for _, sensor := range sensors {
		sensorsBSON = append(sensorsBSON, m.convertSensorToBson(sensor))
	}

	if roundNumber, err = m.getNextRoundCount(ctx); err != nil {
		return err
	}

	round := roundBSON{
		RoundNumber: roundNumber,
		Sensors:     sensorsBSON,
	}

	filter := bson.M{"_id": openMatch.GetID()}
	update := bson.M{
		"$push": bson.M{
			"rounds": round,
		},
	}

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (m *MatchMongo) getCurrentIDCount(ctx context.Context) (int, error) {
	filter := bson.M{"_id": m.collection}

	var result struct {
		ID  string `bson:"_id"`
		Seq int    `bson:"seq"`
	}
	if err := m.client.Database(m.db).Collection("counters").FindOne(ctx, filter).Decode(&result); err != nil {
		return -1, fmt.Errorf("Count not retrieve current ID count: %s", err.Error())
	}

	return result.Seq, nil
}

func (m *MatchMongo) getNextIDCount(ctx context.Context) (int, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)
	filter := bson.M{"_id": m.collection}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	var result struct{ Seq int }
	err := m.client.Database(m.db).Collection("counters").FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return -1, fmt.Errorf("Could not retrieve new id for match: %s", err.Error())
	}

	return result.Seq, nil
}

func (m *MatchMongo) getNextRoundCount(ctx context.Context) (int, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)
	filter := bson.M{"_id": "round"}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	var result struct{ Seq int }
	err := m.client.Database(m.db).Collection("counters").FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return -1, fmt.Errorf("Could not retrieve number for round: %s", err.Error())
	}

	return result.Seq, nil
}

func (m *MatchMongo) resetRoundCount(ctx context.Context) error {
	filter := bson.M{"_id": "round"}
	update := bson.M{"$set": bson.M{"seq": 0}}

	_, err := m.client.Database(m.db).Collection("counters").UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Could not reset number for round: %s", err.Error())
	}

	return nil
}

func (m *MatchMongo) convertReadingToBson(reading entity.Reading) readingBSON {
	return readingBSON{
		Timestamp: reading.GetTimestamp(),
		Value:     reading.GetValue(),
	}
}

func (m *MatchMongo) convertSensorToBson(sensorDTO *dto.SensorDTO) sensorBSON {
	var readings []readingBSON

	for _, reading := range sensorDTO.GetReadings() {
		readings = append(readings, m.convertReadingToBson(reading))
	}

	return sensorBSON{
		Name:     sensorDTO.GetName(),
		Readings: readings,
	}
}

func (m *MatchMongo) convertMatchToBson(matchDTO *dto.MatchDTO) matchBSON {
	return matchBSON{
		ID:           matchDTO.GetID(),
		Title:        matchDTO.GetTitle(),
		Date:         matchDTO.GetDate().Format("2006-01-02"),
		OpponentName: matchDTO.GetOpponentName(),
		Closed:       matchDTO.IsClosed(),
	}
}

func (m *MatchMongo) convertMatchBsonToDto(match matchBSON) (*dto.MatchDTO, error) {
	date, err := time.Parse("2006-01-02", match.Date)
	if err != nil {
		return nil, fmt.Errorf("Could not convert match date string to time: %s", err.Error())
	}
	return dto.NewMatchDTO(match.ID, match.Title, date, match.OpponentName, match.Closed), nil
}
