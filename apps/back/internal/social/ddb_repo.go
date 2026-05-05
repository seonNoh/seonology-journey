package social

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	journeyv1 "github.com/seonNoh/seonology-journey/proto/gen/go/journey/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/seonNoh/seonology-journey/apps/back/internal/repository/ddb"
)

// ─── Companion DDB ───────────────────────────────────────────────────────────

type companionItem struct {
	PK          string `dynamodbav:"PK"` // TRIP#<tripID>
	SK          string `dynamodbav:"SK"` // COMPANION#<memberID>
	TripID      string `dynamodbav:"tripId"`
	MemberID    string `dynamodbav:"memberId"`
	DisplayName string `dynamodbav:"displayName"`
	AvatarURL   string `dynamodbav:"avatarUrl,omitempty"`
	Role        int32  `dynamodbav:"role"`
	InvitedAt   string `dynamodbav:"invitedAt,omitempty"`
}

type CompanionDDBRepo struct {
	client *dynamodb.Client
	table  string
}

func NewCompanionDDBRepo(client *dynamodb.Client) *CompanionDDBRepo {
	return &CompanionDDBRepo{client: client, table: ddb.TableName("companions")}
}

func (r *CompanionDDBRepo) Save(ctx context.Context, c *journeyv1.Companion) error {
	it := companionItem{
		PK: "TRIP#" + c.GetTripId(), SK: "COMPANION#" + c.GetMemberId(),
		TripID: c.GetTripId(), MemberID: c.GetMemberId(),
		DisplayName: c.GetDisplayName(), AvatarURL: c.GetAvatarUrl(),
		Role: int32(c.GetRole()),
	}
	if t := c.GetInvitedAt(); t != nil {
		it.InvitedAt = t.AsTime().Format(time.RFC3339)
	}
	av, err := attributevalue.MarshalMap(it)
	if err != nil {
		return fmt.Errorf("companion marshal: %w", err)
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{TableName: aws.String(r.table), Item: av})
	return err
}

func (r *CompanionDDBRepo) Get(ctx context.Context, tripID, memberID string) (*journeyv1.Companion, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "TRIP#" + tripID},
			"SK": &types.AttributeValueMemberS{Value: "COMPANION#" + memberID},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, ErrNotFound
	}
	var it companionItem
	if err := attributevalue.UnmarshalMap(out.Item, &it); err != nil {
		return nil, err
	}
	return companionFromItem(it), nil
}

func (r *CompanionDDBRepo) ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.Companion, error) {
	keyCond := expression.Key("PK").Equal(expression.Value("TRIP#" + tripID)).
		And(expression.Key("SK").BeginsWith("COMPANION#"))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}
	results := make([]*journeyv1.Companion, 0, len(out.Items))
	for _, item := range out.Items {
		var it companionItem
		if err := attributevalue.UnmarshalMap(item, &it); err != nil {
			continue
		}
		results = append(results, companionFromItem(it))
	}
	return results, nil
}

func (r *CompanionDDBRepo) Delete(ctx context.Context, tripID, memberID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "TRIP#" + tripID},
			"SK": &types.AttributeValueMemberS{Value: "COMPANION#" + memberID},
		},
	})
	return err
}

func (r *CompanionDDBRepo) DeleteByTrip(ctx context.Context, tripID string) error {
	companions, err := r.ListByTrip(ctx, tripID)
	if err != nil {
		return err
	}
	for _, c := range companions {
		if err := r.Delete(ctx, tripID, c.GetMemberId()); err != nil {
			return err
		}
	}
	return nil
}

func companionFromItem(it companionItem) *journeyv1.Companion {
	c := &journeyv1.Companion{
		TripId: it.TripID, MemberId: it.MemberID,
		DisplayName: it.DisplayName, AvatarUrl: it.AvatarURL,
		Role: journeyv1.CompanionRole(it.Role),
	}
	if it.InvitedAt != "" {
		if t, err := time.Parse(time.RFC3339, it.InvitedAt); err == nil {
			c.InvitedAt = timestamppb.New(t)
		}
	}
	return c
}

// ─── Tag DDB ─────────────────────────────────────────────────────────────────

type tagItem struct {
	PK     string `dynamodbav:"PK"` // USER#<userID>
	SK     string `dynamodbav:"SK"` // TAG#<id>
	ID     string `dynamodbav:"id"`
	UserID string `dynamodbav:"userId"`
	Name   string `dynamodbav:"name"`
	Color  string `dynamodbav:"color,omitempty"`
}

type tagTripItem struct {
	PK    string `dynamodbav:"PK"` // TRIP#<tripID>
	SK    string `dynamodbav:"SK"` // TAG#<tagID>
	TagID string `dynamodbav:"tagId"`
}

type TagDDBRepo struct {
	client   *dynamodb.Client
	table    string
	relTable string
}

func NewTagDDBRepo(client *dynamodb.Client) *TagDDBRepo {
	return &TagDDBRepo{client: client, table: ddb.TableName("tags"), relTable: ddb.TableName("tag-trips")}
}

func (r *TagDDBRepo) Save(ctx context.Context, t *journeyv1.Tag) error {
	it := tagItem{
		PK: "USER#" + t.GetUserId(), SK: "TAG#" + t.GetId(),
		ID: t.GetId(), UserID: t.GetUserId(), Name: t.GetName(), Color: t.GetColor(),
	}
	av, err := attributevalue.MarshalMap(it)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{TableName: aws.String(r.table), Item: av})
	return err
}

func (r *TagDDBRepo) Get(ctx context.Context, id string) (*journeyv1.Tag, error) {
	keyCond := expression.Key("GSI1PK").Equal(expression.Value("TAG#" + id))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		IndexName:                 aws.String("GSI1"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1),
	})
	if err != nil {
		return nil, err
	}
	if len(out.Items) == 0 {
		return nil, ErrNotFound
	}
	var it tagItem
	if err := attributevalue.UnmarshalMap(out.Items[0], &it); err != nil {
		return nil, err
	}
	return &journeyv1.Tag{Id: it.ID, UserId: it.UserID, Name: it.Name, Color: it.Color}, nil
}

func (r *TagDDBRepo) ListByUser(ctx context.Context, userID string) ([]*journeyv1.Tag, error) {
	keyCond := expression.Key("PK").Equal(expression.Value("USER#" + userID)).
		And(expression.Key("SK").BeginsWith("TAG#"))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}
	results := make([]*journeyv1.Tag, 0, len(out.Items))
	for _, item := range out.Items {
		var it tagItem
		if err := attributevalue.UnmarshalMap(item, &it); err != nil {
			continue
		}
		results = append(results, &journeyv1.Tag{Id: it.ID, UserId: it.UserID, Name: it.Name, Color: it.Color})
	}
	return results, nil
}

func (r *TagDDBRepo) Delete(ctx context.Context, id string) error {
	t, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "USER#" + t.GetUserId()},
			"SK": &types.AttributeValueMemberS{Value: "TAG#" + id},
		},
	})
	return err
}

func (r *TagDDBRepo) Attach(ctx context.Context, tripID, tagID string) error {
	it := tagTripItem{PK: "TRIP#" + tripID, SK: "TAG#" + tagID, TagID: tagID}
	av, err := attributevalue.MarshalMap(it)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{TableName: aws.String(r.relTable), Item: av})
	return err
}

func (r *TagDDBRepo) Detach(ctx context.Context, tripID, tagID string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.relTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "TRIP#" + tripID},
			"SK": &types.AttributeValueMemberS{Value: "TAG#" + tagID},
		},
	})
	return err
}

func (r *TagDDBRepo) ListByTrip(ctx context.Context, tripID string) ([]*journeyv1.Tag, error) {
	keyCond := expression.Key("PK").Equal(expression.Value("TRIP#" + tripID)).
		And(expression.Key("SK").BeginsWith("TAG#"))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.relTable),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}
	results := make([]*journeyv1.Tag, 0, len(out.Items))
	for _, item := range out.Items {
		var rel tagTripItem
		if err := attributevalue.UnmarshalMap(item, &rel); err != nil {
			continue
		}
		tag, err := r.Get(ctx, rel.TagID)
		if err != nil {
			continue
		}
		results = append(results, tag)
	}
	return results, nil
}

func (r *TagDDBRepo) DetachAllFromTrip(ctx context.Context, tripID string) error {
	keyCond := expression.Key("PK").Equal(expression.Value("TRIP#" + tripID)).
		And(expression.Key("SK").BeginsWith("TAG#"))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.relTable),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return err
	}
	for _, item := range out.Items {
		var rel tagTripItem
		if err := attributevalue.UnmarshalMap(item, &rel); err != nil {
			continue
		}
		_ = r.Detach(ctx, tripID, rel.TagID)
	}
	return nil
}

// ─── Template DDB ────────────────────────────────────────────────────────────

type templateItem struct {
	PK           string `dynamodbav:"PK"` // USER#<userID>
	SK           string `dynamodbav:"SK"` // TEMPLATE#<id>
	ID           string `dynamodbav:"id"`
	UserID       string `dynamodbav:"userId"`
	Name         string `dynamodbav:"name"`
	SourceTripID string `dynamodbav:"sourceTripId,omitempty"`
	CreatedAt    string `dynamodbav:"createdAt"`
	UpdatedAt    string `dynamodbav:"updatedAt"`
}

type TemplateDDBRepo struct {
	client *dynamodb.Client
	table  string
}

func NewTemplateDDBRepo(client *dynamodb.Client) *TemplateDDBRepo {
	return &TemplateDDBRepo{client: client, table: ddb.TableName("templates")}
}

func (r *TemplateDDBRepo) Save(ctx context.Context, t *journeyv1.Template) error {
	it := templateItem{
		PK: "USER#" + t.GetUserId(), SK: "TEMPLATE#" + t.GetId(),
		ID: t.GetId(), UserID: t.GetUserId(), Name: t.GetName(),
		SourceTripID: t.GetSourceTripId(),
	}
	if a := t.GetAudit(); a != nil {
		it.CreatedAt = a.GetCreatedAt().AsTime().Format(time.RFC3339)
		it.UpdatedAt = a.GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	av, err := attributevalue.MarshalMap(it)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{TableName: aws.String(r.table), Item: av})
	return err
}

func (r *TemplateDDBRepo) Get(ctx context.Context, id string) (*journeyv1.Template, error) {
	keyCond := expression.Key("GSI1PK").Equal(expression.Value("TEMPLATE#" + id))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		IndexName:                 aws.String("GSI1"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1),
	})
	if err != nil {
		return nil, err
	}
	if len(out.Items) == 0 {
		return nil, ErrNotFound
	}
	var it templateItem
	if err := attributevalue.UnmarshalMap(out.Items[0], &it); err != nil {
		return nil, err
	}
	return templateFromItem(it), nil
}

func (r *TemplateDDBRepo) ListByUser(ctx context.Context, userID string) ([]*journeyv1.Template, error) {
	keyCond := expression.Key("PK").Equal(expression.Value("USER#" + userID)).
		And(expression.Key("SK").BeginsWith("TEMPLATE#"))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}
	results := make([]*journeyv1.Template, 0, len(out.Items))
	for _, item := range out.Items {
		var it templateItem
		if err := attributevalue.UnmarshalMap(item, &it); err != nil {
			continue
		}
		results = append(results, templateFromItem(it))
	}
	return results, nil
}

func (r *TemplateDDBRepo) Delete(ctx context.Context, id string) error {
	t, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "USER#" + t.GetUserId()},
			"SK": &types.AttributeValueMemberS{Value: "TEMPLATE#" + id},
		},
	})
	return err
}

func templateFromItem(it templateItem) *journeyv1.Template {
	t := &journeyv1.Template{
		Id: it.ID, UserId: it.UserID, Name: it.Name, SourceTripId: it.SourceTripID,
	}
	ca, _ := time.Parse(time.RFC3339, it.CreatedAt)
	ua, _ := time.Parse(time.RFC3339, it.UpdatedAt)
	t.Audit = &journeyv1.AuditTimestamps{CreatedAt: timestamppb.New(ca), UpdatedAt: timestamppb.New(ua)}
	return t
}

// ─── Favorite DDB ────────────────────────────────────────────────────────────

type favoriteItem struct {
	PK            string  `dynamodbav:"PK"` // USER#<userID>
	SK            string  `dynamodbav:"SK"` // FAVORITE#<id>
	ID            string  `dynamodbav:"id"`
	UserID        string  `dynamodbav:"userId"`
	Name          string  `dynamodbav:"name"`
	Lat           float64 `dynamodbav:"lat,omitempty"`
	Lng           float64 `dynamodbav:"lng,omitempty"`
	GooglePlaceID string  `dynamodbav:"googlePlaceId,omitempty"`
	Memo          string  `dynamodbav:"memo,omitempty"`
	CreatedAt     string  `dynamodbav:"createdAt"`
	UpdatedAt     string  `dynamodbav:"updatedAt"`
}

type FavoriteDDBRepo struct {
	client *dynamodb.Client
	table  string
}

func NewFavoriteDDBRepo(client *dynamodb.Client) *FavoriteDDBRepo {
	return &FavoriteDDBRepo{client: client, table: ddb.TableName("favorites")}
}

func (r *FavoriteDDBRepo) Save(ctx context.Context, p *journeyv1.FavoritePlace) error {
	it := favoriteItem{
		PK: "USER#" + p.GetUserId(), SK: "FAVORITE#" + p.GetId(),
		ID: p.GetId(), UserID: p.GetUserId(), Name: p.GetName(),
		GooglePlaceID: p.GetGooglePlaceId(), Memo: p.GetMemo(),
	}
	if loc := p.GetLocation(); loc != nil {
		it.Lat = loc.GetLatitude()
		it.Lng = loc.GetLongitude()
	}
	if a := p.GetAudit(); a != nil {
		it.CreatedAt = a.GetCreatedAt().AsTime().Format(time.RFC3339)
		it.UpdatedAt = a.GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	av, err := attributevalue.MarshalMap(it)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{TableName: aws.String(r.table), Item: av})
	return err
}

func (r *FavoriteDDBRepo) ListByUser(ctx context.Context, userID string) ([]*journeyv1.FavoritePlace, error) {
	keyCond := expression.Key("PK").Equal(expression.Value("USER#" + userID)).
		And(expression.Key("SK").BeginsWith("FAVORITE#"))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}
	results := make([]*journeyv1.FavoritePlace, 0, len(out.Items))
	for _, item := range out.Items {
		var it favoriteItem
		if err := attributevalue.UnmarshalMap(item, &it); err != nil {
			continue
		}
		results = append(results, favoriteFromItem(it))
	}
	return results, nil
}

func (r *FavoriteDDBRepo) Delete(ctx context.Context, id string) error {
	// Need userID to construct key; use GSI1 lookup
	keyCond := expression.Key("GSI1PK").Equal(expression.Value("FAVORITE#" + id))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(r.table),
		IndexName:                 aws.String("GSI1"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1),
	})
	if err != nil {
		return err
	}
	if len(out.Items) == 0 {
		return ErrNotFound
	}
	var it favoriteItem
	if err := attributevalue.UnmarshalMap(out.Items[0], &it); err != nil {
		return err
	}
	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "USER#" + it.UserID},
			"SK": &types.AttributeValueMemberS{Value: "FAVORITE#" + id},
		},
	})
	return err
}

func favoriteFromItem(it favoriteItem) *journeyv1.FavoritePlace {
	p := &journeyv1.FavoritePlace{
		Id: it.ID, UserId: it.UserID, Name: it.Name,
		GooglePlaceId: it.GooglePlaceID, Memo: it.Memo,
	}
	if it.Lat != 0 || it.Lng != 0 {
		p.Location = &journeyv1.GeoPoint{Latitude: it.Lat, Longitude: it.Lng}
	}
	ca, _ := time.Parse(time.RFC3339, it.CreatedAt)
	ua, _ := time.Parse(time.RFC3339, it.UpdatedAt)
	p.Audit = &journeyv1.AuditTimestamps{CreatedAt: timestamppb.New(ca), UpdatedAt: timestamppb.New(ua)}
	return p
}

// ─── Share DDB ───────────────────────────────────────────────────────────────

type shareItem struct {
	PK         string `dynamodbav:"PK"` // SHARE#<code>
	SK         string `dynamodbav:"SK"` // SHARE#<code>
	Code       string `dynamodbav:"code"`
	TripID     string `dynamodbav:"tripId"`
	Permission int32  `dynamodbav:"permission"`
	ExpiresAt  string `dynamodbav:"expiresAt,omitempty"`
	CreatedAt  string `dynamodbav:"createdAt"`
	UpdatedAt  string `dynamodbav:"updatedAt"`
}

type ShareDDBRepo struct {
	client *dynamodb.Client
	table  string
}

func NewShareDDBRepo(client *dynamodb.Client) *ShareDDBRepo {
	return &ShareDDBRepo{client: client, table: ddb.TableName("shares")}
}

func (r *ShareDDBRepo) Save(ctx context.Context, s *journeyv1.Share) error {
	it := shareItem{
		PK: "SHARE#" + s.GetCode(), SK: "SHARE#" + s.GetCode(),
		Code: s.GetCode(), TripID: s.GetTripId(),
		Permission: int32(s.GetPermission()),
	}
	if t := s.GetExpiresAt(); t != nil {
		it.ExpiresAt = t.AsTime().Format(time.RFC3339)
	}
	if a := s.GetAudit(); a != nil {
		it.CreatedAt = a.GetCreatedAt().AsTime().Format(time.RFC3339)
		it.UpdatedAt = a.GetUpdatedAt().AsTime().Format(time.RFC3339)
	}
	av, err := attributevalue.MarshalMap(it)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{TableName: aws.String(r.table), Item: av})
	return err
}

func (r *ShareDDBRepo) Get(ctx context.Context, code string) (*journeyv1.Share, error) {
	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "SHARE#" + code},
			"SK": &types.AttributeValueMemberS{Value: "SHARE#" + code},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, ErrNotFound
	}
	var it shareItem
	if err := attributevalue.UnmarshalMap(out.Item, &it); err != nil {
		return nil, err
	}
	return shareFromItem(it), nil
}

func (r *ShareDDBRepo) Delete(ctx context.Context, code string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "SHARE#" + code},
			"SK": &types.AttributeValueMemberS{Value: "SHARE#" + code},
		},
	})
	return err
}

func shareFromItem(it shareItem) *journeyv1.Share {
	s := &journeyv1.Share{
		Code: it.Code, TripId: it.TripID,
		Permission: journeyv1.CompanionRole(it.Permission),
	}
	if it.ExpiresAt != "" {
		if t, err := time.Parse(time.RFC3339, it.ExpiresAt); err == nil {
			s.ExpiresAt = timestamppb.New(t)
		}
	}
	ca, _ := time.Parse(time.RFC3339, it.CreatedAt)
	ua, _ := time.Parse(time.RFC3339, it.UpdatedAt)
	s.Audit = &journeyv1.AuditTimestamps{CreatedAt: timestamppb.New(ca), UpdatedAt: timestamppb.New(ua)}
	return s
}
