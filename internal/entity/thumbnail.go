package entity

type Thumbnail struct {
	ID        int64  `db:"id"`
	VideoURL  string `db:"video_url"`
	Thumbnail []byte `db:"grpc"`
}
