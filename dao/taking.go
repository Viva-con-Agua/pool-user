package dao

import (
	"context"
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago/vmdb"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func TakingInsert(ctx context.Context, i *models.TakingCreate, token *vcapool.AccessToken) (r *models.Taking, err error) {
	//create taking model form i.
	taking := i.TakingDatabase()
	if err = TakingCollection.InsertOne(ctx, taking); err != nil {
		return
	}
	//create sources
	for _, source := range i.NewSource {
		if source.HasExternal {
			source.External.ReasonForPayment, _ = GetNewReasonForPayment(ctx, i.CrewID)
			deposit := &models.DepositDatabase{
				ID:               uuid.NewString(),
				ReasonForPayment: source.External.ReasonForPayment,
				Status:           "wait",
				Money:            source.Money,
			}
			depositUnit := &models.DepositUnit{
				ID:        uuid.NewString(),
				TakingID:  taking.ID,
				Money:     source.Money,
				DepositID: deposit.ID,
				Status:    "wait",
			}
			if err = DepositCollection.InsertOne(ctx, deposit); err != nil {
				return
			}
			if err = DepositUnitCollection.InsertOne(ctx, depositUnit); err != nil {
				return
			}

		}
	}
	if i.NewSource != nil {
		sources := i.SourceList(taking.ID)
		if err = SourceCollection.InsertMany(ctx, sources.InsertMany()); err != nil {
			return
		}
	}
	r = new(models.Taking)
	if err = TakingCollection.AggregateOne(
		ctx,
		models.NewTakingsPipeline().Match(bson.D{{Key: "_id", Value: taking.ID}}).Pipe,
		r,
	); err != nil {
		return
	}
	return
}

func TakingUpdate(ctx context.Context, i *models.TakingUpdate) (r *models.Taking, err error) {
	takingDatabase := new(models.TakingDatabase)
	if err = TakingCollection.FindOne(ctx, bson.D{{Key: "_id", Value: i.ID}}, takingDatabase); err != nil {
		return
	}
	i.State = &takingDatabase.State
	for _, v := range i.Sources {
		//create new sources
		if v.ID == "" {
			i.State.Open.Amount += v.Money.Amount
			v.TakingID = i.ID
			newSource := v.Source()
			if err = SourceCollection.InsertOne(ctx, newSource); err != nil {
				return
			}
		}
		if v.UpdateState == "deleted" {
			deleteSource := new(models.Source)
			if err = SourceCollection.FindOne(ctx, bson.D{{Key: "_id", Value: v.ID}}, deleteSource); err != nil {
				return
			}
			takingDatabase.State.Open.Amount -= deleteSource.Money.Amount
			if err = SourceCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: deleteSource.ID}}); err != nil {
				return
			}
		}
		if v.UpdateState == "updated" {
			databaseSource := new(models.Source)
			if err = SourceCollection.FindOne(
				ctx,
				bson.D{{Key: "_id", Value: v.ID}},
				databaseSource,
			); err != nil {
				return
			}
			if v.Money.Amount != databaseSource.Money.Amount {
				i.State.Open.Amount -= databaseSource.Money.Amount
				i.State.Open.Amount += v.Money.Amount
			}
			if err = SourceCollection.UpdateOne(
				ctx,
				bson.D{{Key: "_id", Value: v.ID}},
				vmdb.UpdateSet(v),
				nil,
			); err != nil {
				return
			}
		}
	}
	r = new(models.Taking)
	if err = TakingCollection.UpdateOneAggregate(
		ctx,
		bson.D{{Key: "_id", Value: i.ID}},
		vmdb.UpdateSet(i),
		r,
		models.NewTakingsPipeline().Match(bson.D{{Key: "_id", Value: i.ID}}).Pipe,
	); err != nil {
		return
	}
	return
}
